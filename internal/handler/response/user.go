package response

type UserInfoResponse struct {
	UserId    string `json:"user_id"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Telephone string `json:"telephone"`
}
