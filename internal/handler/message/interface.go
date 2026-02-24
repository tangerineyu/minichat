package message

import "github.com/gin-gonic/gin"

type MessageHandlerInterface interface {
	// SendMessage 发送消息
	SendMessage(c *gin.Context)
	// GetMessageHistory 获取消息历史记录
	GetMessageHistory(c *gin.Context)
	// WithdrawMessage 撤回消息
	WithdrawMessage(c *gin.Context)
}
