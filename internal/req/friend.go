package req

type BlackFriendReq struct {
	FriendId int64 `json:"friend_id" binding:"required"`
}

type FriendListReq struct {
	FriendId int64  `json:"friend_id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}
