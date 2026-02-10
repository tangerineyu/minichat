package dto

import "time"

// FriendApplyItem 用于对外返回好友申请列表条目。
type FriendApplyItem struct {
	ApplyId  int64  `json:"apply_id"`
	ApplyMsg string `json:"apply_msg"`
	Status   int8   `json:"status"`

	FromId           int64     `json:"from_user_id"`
	FromUserNickname string    `json:"from_user_nickname"`
	FromUserAvatar   string    `json:"from_user_avatar"`
	CreatedAt        time.Time `json:"created_at"`
}
