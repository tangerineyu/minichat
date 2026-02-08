package user

import (
	"context"
	"minichat/internal/handler/response"
	"minichat/internal/req"
)

type UserServiceInterface interface {
	// 注册
	Register(ctx context.Context, in req.RegisterReq) error
	// 登录，返回 access token 和 refresh token
	Login(ctx context.Context, in req.LoginReq) (string, string, error)
	// 更改基础信息，昵称，头像
	UpdateUserInfo(ctx context.Context, userId int64, in req.UpdateUserReq) error
	// 更改密码
	ChangePassword(ctx context.Context, id int64, in req.ChangePasswordReq) error
	// 更改用户Id
	ChangeUserId(ctx context.Context, id int64, in req.ChangeUserIdReq) error
	// 获取用户信息
	GetUserInfo(ctx context.Context, id int64) (response.UserInfoResponse, error)
	// 用户注销（需要二次密码校验）
	CancelAccount(ctx context.Context, id int64, password string) error
}
