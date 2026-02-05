package user

import (
	"context"
	"errors"
	"fmt"
	commonModel "minichat/internal/common"
	"minichat/internal/model"
	repo "minichat/internal/repo/user"
	"minichat/internal/req"
	jwtutil "minichat/util/jwt"
	crypto "minichat/util/password"
	"minichat/util/snowflake"
	"time"
)

type UserService struct {
	userRepo repo.UserRepoInterface
}

// 用户注册
func (u UserService) Register(ctx context.Context, in req.RegisterReq) error {
	if in.Telephone == "" || len(in.Password) < 6 {
		return errors.New("参数错误")
	}

	exists, err := u.userRepo.CheckUserExistsByTelephone(ctx, in.Telephone)
	if err != nil {
		return errors.New("查询用户失败：" + err.Error())
	}
	if exists {
		return errors.New(commonModel.TELEPHONE_HAS_REGISTERED)
	}
	pwd, err := crypto.HashPassword(in.Password)
	if err != nil {
		return err
	}
	//
	username := snowflake.GenStringID()
	newUser := &model.User{
		Username:  username,
		Telephone: in.Telephone,
		Password:  pwd,
		Nickname:  in.Nickname,
		Avatar:    "default_avatar.png",
		Status:    1,
	}
	if err := u.userRepo.CreateUser(ctx, newUser); err != nil {
		return errors.New("创建用户失败：" + err.Error())
	}
	return nil
}

// 用户登录
func (u UserService) Login(ctx context.Context, in req.LoginReq) (string, string, error) {
	if len(in.Password) < 6 {
		return "", "", errors.New(commonModel.PASSWORD_LENGTH_LESS_SIX)
	}
	if in.Username == "" {
		return "", "", errors.New(commonModel.USERNAME_IS_EMPTY)
	}
	user, err := u.userRepo.GetUserByUsername(ctx, in.Username)
	if err != nil {
		return "", "", errors.New("查询用户失败：" + err.Error())
	}
	if !crypto.CheckPasswordHash(in.Password, user.Password) {
		return "", "", errors.New(commonModel.PASSWORD_ERROR)
	}

	access, refresh, err := jwtutil.GenerateToken(*user)
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

func (u UserService) UpdateUserInfo(ctx context.Context, userId int64, in req.UpdateUserReq) error {
	user, err := u.userRepo.GetUserById(ctx, userId)
	if err != nil {
		return errors.New("查询用户失败：" + err.Error())
	}
	user.Nickname = in.Nickname
	user.Avatar = in.Avatar

	if err := u.userRepo.UpdateUser(ctx, user); err != nil {
		return errors.New("更新用户信息失败：" + err.Error())
	}
	return nil
}

func (u UserService) ChangePassword(ctx context.Context, id int64, in req.ChangePasswordReq) error {
	user, err := u.userRepo.GetUserById(ctx, id)
	if err != nil {
		return errors.New("查询用户失败：" + err.Error())
	}
	// 检查旧密码是否正确
	if !crypto.CheckPasswordHash(in.OldPassword, user.Password) {
		return errors.New(commonModel.PASSWORD_ERROR)
	}
	if in.NewPassword != in.OldPassword {
		return errors.New("新密码不能与旧密码相同")
	}
	newPasswordHash, err := crypto.HashPassword(in.NewPassword)
	if err != nil {
		return errors.New("密码加密失败：" + err.Error())
	}
	return u.userRepo.ChangePassword(ctx, user.ID, newPasswordHash)
}

const UsernameChangeIntervalDays = 180 * 24 * time.Hour

func (u UserService) ChangeUsername(ctx context.Context, id int64, in req.ChangeUsernameReq) error {
	user, err := u.userRepo.GetUserById(ctx, id)
	if err != nil {
		return errors.New("查询用户失败：" + err.Error())
	}
	// 检查上次更改过了多长时间
	if user.UsernameChangedAt != nil {
		gap := time.Since(*user.UsernameChangedAt)
		if gap < UsernameChangeIntervalDays {
			remainingDays := int((UsernameChangeIntervalDays - gap).Hours() / 24)
			return fmt.Errorf("短期内只能更改一次，请在 %d 天后再试", remainingDays)
		}
	}
	// 检查新用户名是否已存在
	existsUser, err := u.userRepo.GetUserByUsername(ctx, in.NewUsername)
	if err == nil && existsUser.ID != user.ID {
		return errors.New(commonModel.USERNAME_HAS_EXISTS)
	}
	now := time.Now()
	user.UsernameChangedAt = &now
	return u.userRepo.ChangeUsername(ctx, in.NewUsername, user.UsernameChangedAt)
}

func NewUserService(userRepo repo.UserRepoInterface) UserServiceInterface {
	return &UserService{
		userRepo: userRepo,
	}
}
