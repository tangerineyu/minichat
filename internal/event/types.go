package event

// WsTargetType 表示 websocket 推送目标类型。
// 默认（空字符串）会按 User 处理，以兼容旧逻辑。
type WsTargetType string

const (
	WsTargetUser  WsTargetType = "user"
	WsTargetGroup WsTargetType = "group"
)

// WsEnvelope 是事件总线里承载 ws 推送的数据封装。

// 群聊/群通知：填 TargetType=group 且 TargetID=groupID，并把 GroupMemberIDs 填好（推荐）。
// 支持直接填 TargetUserIDs 做批量推送。
type WsEnvelope struct {
	// TargetType 决定 TargetID 的含义（userID / groupID）。空值默认按 user。
	TargetType WsTargetType `json:"target_type"`

	// TargetID 旧字段：
	//  - TargetType=user  -> userID
	//  - TargetType=group -> groupID
	TargetID int64 `json:"target_id"`

	// TargetUserIDs 可选：用于一次事件推送到多个用户（比如群聊 fan-out 后的成员列表）。
	TargetUserIDs []int64 `json:"target_user_ids"`

	// GroupMemberIDs 可选：当 TargetType=group 时，建议把成员 ID 列表放这里。
	//（与 TargetUserIDs 二选一即可）
	GroupMemberIDs []int64 `json:"group_member_ids"`

	Event string `json:"event"`
	Data  any    `json:"data"`
}

type ChatMessagePayload struct {
	MsgID      int64  `json:"msg_id,string"`
	SenderID   int64  `json:"sender_id,string"`
	ReceiverID int64  `json:"receiver_id,string"`
	Content    string `json:"content"`
	MsgType    int8   `json:"msg_type"`
}

type ChatWithdrawPayload struct {
	MsgID      int64 `json:"msg_id,string"`
	WithdrawBy int64 `json:"withdraw_by,string"`
}

type FriendApplyPayload struct {
	ApplyID   int64  `json:"apply_id,string"`
	Applicant int64  `json:"applicant,string"` // 申请人
	Msg       string `json:"msg"`
}
