package message

import (
	"context"
	"minichat/internal/model"

	"gorm.io/gorm"
)

var _ MessageRepoInterface = (*MessageRepo)(nil)

type MessageRepo struct {
	db *gorm.DB
}

func (m *MessageRepo) WithdrawMessage(ctx context.Context, msgId int64) error {
	return m.db.WithContext(ctx).Model(model.Message{}).
		Where("id = ?", msgId).
		Updates(map[string]interface{}{
			"status":  1,  // 1 表示已撤回
			"content": "", //清空内容
		}).Error
}

func (m *MessageRepo) GetMessageById(ctx context.Context, msgId int64) (*model.Message, error) {
	var msg *model.Message
	err := m.db.WithContext(ctx).Where("id=?", msgId).First(&msg).Error
	return msg, err
}

func (m *MessageRepo) SendMessage(ctx context.Context, senderId int64, msg *model.Message) error {
	// 兜底确保 senderId 一致
	if msg != nil {
		msg.SenderID = senderId
	}
	return m.db.WithContext(ctx).Create(msg).Error
}

func (m *MessageRepo) GetMessageList(ctx context.Context, userId int64, targetId int64, sessionType int8, beforeID int64, limit int) ([]*model.Message, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	var msgs []*model.Message
	query := m.db.WithContext(ctx).Model(&model.Message{})

	switch sessionType {
	case 1:
		// 私聊: (sender=Me AND receiver=Peer) OR (sender=Peer AND receiver=Me)
		query = query.Where("((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)) AND session_type = 1", userId, targetId, targetId, userId)
	case 2:
		// 群聊: receiver_id = GroupID AND session_type = 2
		query = query.Where("receiver_id = ? AND session_type = 2", targetId)
	default:
		return nil, nil
	}

	if beforeID > 0 {
		query = query.Where("id < ?", beforeID)
	}

	err := query.
		Order("id desc").
		Limit(limit).
		Find(&msgs).Error
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func NewMessageRepo(db *gorm.DB) MessageRepoInterface {
	return &MessageRepo{db: db}
}
