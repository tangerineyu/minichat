package dto

// FriendItem 用于对外返回好友列表条目。
type FriendItem struct {
	FriendId int64  `json:"friend_id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}
