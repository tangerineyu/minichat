package message

import (
	"minichat/internal/handler/response"
	"minichat/internal/req"
	messageService "minichat/internal/service/message"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	messageService messageService.MessageServiceInterface
}

func (m *MessageHandler) WithdrawMessage(c *gin.Context) {
	var in req.WithDrawMsgReq
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, 400, "invalid request")
		return
	}
	id := c.GetInt64("id")
	if err := m.messageService.WithdrawMessage(c.Request.Context(), id, in); err != nil {
		response.ServerError(c, err)
		return
	}
	response.Success(c, "撤回成功")
}

func (m *MessageHandler) GetMessageHistory(c *gin.Context) {
	userId := c.GetInt64("id")
	var in req.GetMessageHistoryReq
	if err := c.ShouldBindQuery(&in); err != nil {
		response.Fail(c, 400, "invalid request")
		return
	}
	msgs, err := m.messageService.GetMessageHistory(c.Request.Context(), userId, in)
	if err != nil {
		response.ServerError(c, err)
		return
	}

	// next_cursor: 本页最后一条消息的 id（因为是 id desc）
	nextCursor := ""
	if len(msgs) > 0 {
		nextCursor = strconv.FormatInt(msgs[len(msgs)-1].ID, 10)
	}
	response.Success(c, gin.H{"list": msgs, "next_cursor": nextCursor})
}

func (m *MessageHandler) SendMessage(c *gin.Context) {
	var in req.SendMessageReq
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}

	senderId := c.GetInt64("id")
	msg, err := m.messageService.SendMessage(
		c.Request.Context(),
		senderId,
		in.ReceiverId,
		in.SessionType,
		in.MsgType,
		in.Content,
	)
	if err != nil {
		response.ServerError(c, err)
		return
	}
	response.Success(c, gin.H{"message": msg})
}

func NewMessageHandler(messageService messageService.MessageServiceInterface) MessageHandlerInterface {
	return &MessageHandler{messageService: messageService}
}
