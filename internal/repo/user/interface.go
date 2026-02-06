package user

import (
	"context"
	"minichat/internal/model"
	"time"
)

type UserRepoInterface interface {
	// CreateUser 创建一个新的用户
	CreateUser(ctx context.Context, newUser *model.User) error
	// GetUserByUserId 通过用户名查询用户
	GetUserByUserId(ctx context.Context, userId string) (*model.User, error)
	// GetUserByTelephone 通过手机号查询用户
	GetUserByTelephone(ctx context.Context, telephone string) (*model.User, error)
	// GetUserById 通过ID查询用户
	GetUserById(ctx context.Context, id int64) (*model.User, error)
	// CheckUserExistsByTelephone 检查手机号是否存在
	CheckUserExistsByTelephone(ctx context.Context, telephone string) (bool, error)
	// UpdateUser 更新用户基础信息 (头像、昵称等)
	UpdateUser(ctx context.Context, user *model.User) error
	// ChangePassword 修改用户密码
	ChangePassword(ctx context.Context, userID int64, newHashedPassword string) error
	// UpdateUserId 更新userId
	ChangeUserId(ctx context.Context, id int64, username string, usernameChangedAt *time.Time) error
}
