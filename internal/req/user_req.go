package req

// RegisterReq represents the structure for user registration requests.
type RegisterReq struct {
	//Username  string `json:"username"`
	Telephone string `json:"telephone"`
	Password  string `json:"password"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
}

// Login 支持 username 或 telephone 二选一
// Password 必填
type LoginReq struct {
	UserId    string `json:"user_id"`
	Telephone string `json:"telephone"`
	Password  string `json:"password"`
}

// UpdateUserReq represents the structure for updating user information.
type UpdateUserReq struct {
	Nickname *string `json:"nickname"`
	Avatar   *string `json:"avatar"`
}

// 更改密码
type ChangePasswordReq struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type CancelAccountReq struct {
	Password string `json:"password"`
}
type ChangeUserIdReq struct {
	NewUserId string `json:"user_id"`
}

// 注意：AuthResponse 属于响应 DTO（且不建议直接暴露 model.User），已迁移出 req 包。
