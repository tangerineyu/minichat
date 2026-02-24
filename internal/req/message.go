package req

type SendMessageReq struct {
	ReceiverId  int64  `json:"receiver_id" binding:"required"`
	SessionType int8   `json:"session_type"` // 1私聊, 2群聊
	MsgType     int8   `json:"msg_type"`     // 1文本, 2图片, 3文件, 4系统通知
	Content     string `json:"content" binding:"required"`
}

// GetMessageHistoryReq 获取消息历史（滚动分页）
// cursor: 上一页返回的 next_cursor（message_id），首次请求可不传
// limit: 期望返回条数，默认 20
type GetMessageHistoryReq struct {
	TargetID    int64 `form:"target_id" json:"target_id" binding:"required"`       // 对方ID(私聊) 或 群ID(群聊)
	SessionType int8  `form:"session_type" json:"session_type" binding:"required"` // 1私聊, 2群聊
	CursorPageReq
}

type WithDrawMsgReq struct {
	MsgId       int64 `json:"msg_id" binding:"required"`
	ReceiverId  int64 `json:"receiver_id" binding:"required"`
	SessionType int8  `json:"session_type" binding:"required"`
}
