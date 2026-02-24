package message

import (
	"context"
	"encoding/json"
	"errors"
	"minichat/internal/event"
	"minichat/internal/req"
	"minichat/internal/websocket"
	"strconv"
	"strings"
	"time"

	"minichat/internal/model"
	friendRepo "minichat/internal/repo/friend"
	groupMemberRepo "minichat/internal/repo/group_member"
	messageRepo "minichat/internal/repo/message"
	"minichat/util/snowflake"

	"gorm.io/gorm"
)

var _ MessageServiceInterface = (*MessageService)(nil)

type MessageService struct {
	messageRepo     messageRepo.MessageRepoInterface
	friendRepo      friendRepo.FriendRepoInterface
	groupMemberRepo groupMemberRepo.GroupMemberRepoInterface
}

func (m *MessageService) WithdrawMessage(ctx context.Context, operatorId int64, in req.WithDrawMsgReq) error {
	msg, err := m.messageRepo.GetMessageById(ctx, in.MsgId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("消息不存在")
		}
		return err
	}
	if msg.SenderID != operatorId {
		return errors.New("只能撤回自己的消息")
	}
	// 时间校验
	if time.Since(msg.CreatedAt) > 2*time.Minute {
		return errors.New("发送时间超过2分钟，无法撤回")
	}
	// 更新数据库
	if err := m.messageRepo.WithdrawMessage(ctx, in.MsgId); err != nil {
		return err
	}

	envelope := event.WsEnvelope{
		TargetID: in.ReceiverId,
		Event:    "chat.message.withdraw",
		Data: event.ChatWithdrawPayload{
			MsgID:      in.MsgId,
			WithdrawBy: operatorId,
		},
	}
	event.GlobalBus.Publish(ctx, event.EventMessageWithdraw, envelope)
	return nil
}

func (m *MessageService) GetMessageHistory(ctx context.Context, userId int64, in req.GetMessageHistoryReq) ([]*model.Message, error) {
	if userId <= 0 {
		return nil, errors.New("invalid user_id")
	}
	if in.TargetID <= 0 {
		return nil, errors.New("invalid target_id")
	}

	limit := in.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	var beforeID int64
	if strings.TrimSpace(in.Cursor) != "" {
		id, err := strconv.ParseInt(in.Cursor, 10, 64)
		if err != nil {
			return nil, errors.New("invalid cursor")
		}
		beforeID = id
	}

	return m.messageRepo.GetMessageList(ctx, userId, in.TargetID, in.SessionType, beforeID, limit)
}

func (m *MessageService) SendMessage(ctx context.Context, senderId int64, receiverId int64, sessionType int8, msgType int8, content string) (*model.Message, error) {
	if senderId <= 0 {
		return nil, errors.New("invalid sender")
	}
	if receiverId <= 0 {
		return nil, errors.New("invalid receiver")
	}
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, errors.New("content is empty")
	}
	if sessionType == 0 {
		sessionType = 1
	}
	if msgType == 0 {
		msgType = 1
	}

	// 私聊：必须是好友（且不能在黑名单）
	if sessionType == 1 {
		if m.friendRepo == nil {
			return nil, errors.New("friend repo not configured")
		}
		relation, err := m.friendRepo.GetFriendRelation(ctx, senderId, receiverId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("你们不是好友关系，无法发送消息")
			}
			return nil, err
		}
		if relation == nil || relation.Status != 1 {
			return nil, errors.New("你们不是好友关系，无法发送消息")
		}
	}

	// 群聊：必须是群成员（避免非群成员知道 groupId 后随意发消息）
	if sessionType == 2 {
		if m.groupMemberRepo == nil {
			return nil, errors.New("group member repo not configured")
		}
		gm, err := m.groupMemberRepo.GetMemberById(ctx, senderId, receiverId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("你不是该群成员，无法发送消息")
			}
			return nil, err
		}
		// Status: 0正常/1禁言/2退出
		if gm == nil || gm.Status != 0 {
			return nil, errors.New("你不是该群成员，无法发送消息")
		}
	}

	msg := &model.Message{
		ID:          snowflake.GenInt64ID(),
		SenderID:    senderId,
		ReceiverID:  receiverId,
		SessionType: sessionType,
		MsgType:     msgType,
		Content:     content,
		CreatedAt:   time.Now(),
	}

	if err := m.messageRepo.SendMessage(ctx, senderId, msg); err != nil {
		return nil, err
	}

	// websocket 消息推送：
	//   私聊：推给对方 receiverId（在线则收到，不在线忽略）
	//   群聊：receiverId 表示 groupId，需要推给所有群成员 userId
	pushData, _ := json.Marshal(msg)

	switch sessionType {
	case 1:
		go websocket.GlobalHub.PushMessage(receiverId, pushData)
	case 2:
		// 群聊：查成员列表并逐个推送
		if m.groupMemberRepo != nil {
			members, err := m.groupMemberRepo.GetGroupMembers(ctx, receiverId)
			if err == nil {
				for _, mem := range members {
					if mem == nil {
						continue
					}
					// 是否回推给自己。
					// 这里选择“也推给自己”，方便多端/前端统一用推送来更新界面。
					if mem.Status != 0 {
						continue
					}
					uid := mem.UserID
					go websocket.GlobalHub.PushMessage(uid, pushData)
				}
			}
		}
	}

	return msg, nil
}

func NewMessageService(
	messageRepo messageRepo.MessageRepoInterface,
	friendRepo friendRepo.FriendRepoInterface,
	groupMemberRepo groupMemberRepo.GroupMemberRepoInterface,
) MessageServiceInterface {
	return &MessageService{messageRepo: messageRepo, friendRepo: friendRepo, groupMemberRepo: groupMemberRepo}
}
