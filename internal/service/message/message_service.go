package message

import (
	"context"
	"encoding/json"
	"errors"
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

type MessageService struct {
	messageRepo     messageRepo.MessageRepoInterface
	friendRepo      friendRepo.FriendRepoInterface
	groupMemberRepo groupMemberRepo.GroupMemberRepoInterface
}

var _ MessageServiceInterface = (*MessageService)(nil)

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
	// websocket 消息推送等后续处理可在这里进行（异步）
	pushData, _ := json.Marshal(msg)
	go func(targetId int64, data []byte) {
		websocket.GlobalHub.PushMessage(targetId, data)
	}(receiverId, pushData)
	return msg, nil
}

func NewMessageService(
	messageRepo messageRepo.MessageRepoInterface,
	friendRepo friendRepo.FriendRepoInterface,
	groupMemberRepo groupMemberRepo.GroupMemberRepoInterface,
) MessageServiceInterface {
	return &MessageService{messageRepo: messageRepo, friendRepo: friendRepo, groupMemberRepo: groupMemberRepo}
}
