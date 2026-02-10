package dto

import "time"

// FriendItem 用于对外返回好友列表条目。
type FriendItem struct {
	FriendId       int64     `json:"friend_id"`
	FriendNickname string    `json:"friend_nickname"`
	FriendRemark   string    `json:"friend_remark"`
	FriendAvatar   string    `json:"friend_avatar"`
	CreateAt       time.Time `json:"create_at"`
}
