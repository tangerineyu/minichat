package req

type SendFriendApplyReq struct {
	//FromUserId string `json:"from_user_id"`
	ToUserId int64  `json:"to_user_id"`
	Message  string `json:"message"`
}

type DealWithFriendApplyReq struct {
	ApplyId int64 `json:"apply_id"`
	Action  int8  `json:"action" comment:"同意 1，拒绝 0"`
}
