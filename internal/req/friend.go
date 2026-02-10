package req

type BlackFriendReq struct {
	FriendId int64 `json:"friend_id" binding:"required"`
}

type UpdateFriendRemarkReq struct {
	FriendId int64  `json:"friend_id" binding:"required"`
	Remark   string `json:"remark" binding:"required"`
}

// FriendListReq 属于响应 DTO，已迁移到 internal/dto.FriendItem。
