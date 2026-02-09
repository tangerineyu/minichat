package req

import "time"

type SendFriendApplyReq struct {
	//FromUserId string `json:"from_user_id"`
	ToUserId int64  `json:"to_user_id"`
	Message  string `json:"message"`
}

type DealWithFriendApplyReq struct {
	ApplyId int64 `json:"apply_id"`
	Action  int8  `json:"action" comment:"同意 1，拒绝 0"`
	// 可选字段，用于添加好友时的备注，默认是对方的昵称
	Remark string `json:"remark"`
}

type FriendApplyListReq struct {
	ApplyId  int64  `json:"apply_id"`
	ApplyMsg string `json:"apply_msg"`
	Status   int8   `json:"status"`

	FromId           int64     `json:"from_user_id"`
	FromUserNickname string    `json:"from_user_nickname"`
	FromUserAvatar   string    `json:"from_user_avatar"`
	CreatedAt        time.Time `json:"created_at"`
}
