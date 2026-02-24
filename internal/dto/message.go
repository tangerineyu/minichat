package dto

// websocket推送给前端的撤回事件
type WithdrawEvent struct {
	Type       string `json:"type"` // "withdraw"
	MsgId      int64  `json:"msg_id"`
	WithdrawBy int64  `json:"withdraw_by"` // 撤回者ID
}
