package message

import (
	"context"
	"minichat/internal/model"
)

type MessageRepoInterface interface {
	// SendMessage 发送消息
	SendMessage(ctx context.Context, senderId int64, msg *model.Message) error
	// GetMessageList 获取消息列表（滚动分页）
	// beforeID: 为空/0 表示从最新开始；否则取 message.id < beforeID 的更早消息
	GetMessageList(ctx context.Context, userId int64, targetId int64, sessionType int8, beforeID int64, limit int) ([]*model.Message, error)
	// GetMessageById 通过ID获取消息
	GetMessageById(ctx context.Context, msgId int64) (*model.Message, error)
	// WithdrawMessage 撤回消息
	WithdrawMessage(ctx context.Context, msgId int64) error
}
