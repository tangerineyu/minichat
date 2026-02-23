package message

import (
	"context"
	"minichat/internal/model"
	"minichat/internal/req"
)

type MessageServiceInterface interface {
	// SendMessage 发送消息
	SendMessage(ctx context.Context, senderId int64, receiverId int64, sessionType int8, msgType int8, content string) (*model.Message, error)
	// GetMessageHistory 获取消息历史
	GetMessageHistory(ctx context.Context, userId int64, in req.GetMessageHistoryReq) ([]*model.Message, error)
}
