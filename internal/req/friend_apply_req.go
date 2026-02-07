package req

type SendFriendApplyReq struct {
	//FromUserId string `json:"from_user_id"`
	ToUserId string `json:"to_user_id"`
	Message  string `json:"message"`
}
