package req

import "minichat/internal/model"

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
	Username  string `json:"username"`
	Telephone string `json:"telephone"`
	Password  string `json:"password"`
}
type UpdateUserReq struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// 更改密码
type ChangePasswordReq struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
type ChangeUsernameReq struct {
	NewUsername string `json:"username"`
}
type AuthResponse struct {
	Token string     `json:"token"`
	User  model.User `json:"user"`
}
