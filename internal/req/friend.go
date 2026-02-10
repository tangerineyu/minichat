package req

type BlackFriendReq struct {
	FriendId int64 `json:"friend_id" binding:"required"`
}

// FriendListReq 属于响应 DTO，已迁移到 internal/dto.FriendItem。
