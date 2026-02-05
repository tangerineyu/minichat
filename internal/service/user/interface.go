package user

import (
	"context"
	"minichat/internal/req"
)

type UserServiceInterface interface {
	Register(ctx context.Context, in req.RegisterReq) error
	Login(ctx context.Context, in req.LoginReq) (string, string, error)
	UpdateUserInfo(ctx context.Context, userId int64, in req.UpdateUserReq) error

	ChangePassword(ctx context.Context, id int64, in req.ChangePasswordReq) error

	ChangeUsername(ctx context.Context, id int64, in req.ChangeUsernameReq) error
}
