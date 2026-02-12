package user

import (
	"context"
	"errors"
	"fmt"
	"io"
	commonModel "minichat/internal/common"
	"minichat/internal/dto"
	"minichat/internal/model"
	repo "minichat/internal/repo/user"
	"minichat/internal/req"
	jwtutil "minichat/util/jwt"
	crypto "minichat/util/password"
	"minichat/util/snowflake"
	ossutil "minichat/util/storage/oss"
	"path"
	"strings"
	"time"
)

type UserService struct {
	userRepo repo.UserRepoInterface
	oss      ossutil.OSSInterface
}

func (u *UserService) CancelAccount(ctx context.Context, id int64, password string) error {
	if len(password) < 6 {
		return errors.New(commonModel.PASSWORD_LENGTH_LESS_SIX)
	}
	user, err := u.userRepo.GetUserById(ctx, id)
	if err != nil {
		return errors.New("查询用户失败：" + err.Error())
	}
	if user.Status == 2 {
		return errors.New("账号已注销")
	}
	// 校验密码
	if !crypto.CheckPasswordHash(password, user.Password) {
		return errors.New(commonModel.PASSWORD_ERROR)
	}

	user.Status = 2
	// 为了避免手机号被占用，注销账号时在手机号前添加前缀，并加上时间戳
	user.Telephone = fmt.Sprintf("canceled_%d_%s", time.Now().Unix(), user.Telephone)
	// 注销账号时将用户名改为 "已注销用户"，头像改为默认头像
	user.Nickname = "已注销用户"
	user.Avatar = "default_avatar.png"

	return u.userRepo.UpdateUser(ctx, user)
}

// 用户注册
func (u *UserService) Register(ctx context.Context, in req.RegisterReq) error {
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
		UserId:    username,
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
func (u *UserService) Login(ctx context.Context, in req.LoginReq) (string, string, error) {
	if len(in.Password) < 6 {
		return "", "", errors.New(commonModel.PASSWORD_LENGTH_LESS_SIX)
	}
	if in.UserId == "" && in.Telephone == "" {
		return "", "", errors.New(commonModel.USERID_AND_TELEPHONE_IS_EMPTY)
	}
	var user *model.User
	var err error
	if in.Telephone != "" {
		user, err = u.userRepo.GetUserByTelephone(ctx, in.Telephone)
	} else {
		user, err = u.userRepo.GetUserByUserId(ctx, in.UserId)
	}
	//user, err := u.userRepo.GetUserByUserId(ctx, in.UserId)
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

func (u *UserService) UpdateUserInfo(ctx context.Context, userId int64, in req.UpdateUserReq) error {
	user, err := u.userRepo.GetUserById(ctx, userId)
	if err != nil {
		return errors.New("查询用户失败：" + err.Error())
	}
	if in.Nickname != nil {
		user.Nickname = *in.Nickname
	}
	if in.Avatar != nil {
		user.Avatar = *in.Avatar
	}

	if err := u.userRepo.UpdateUser(ctx, user); err != nil {
		return errors.New("更新用户信息失败：" + err.Error())
	}
	return nil
}

func (u *UserService) ChangePassword(ctx context.Context, id int64, in req.ChangePasswordReq) error {
	user, err := u.userRepo.GetUserById(ctx, id)
	if err != nil {
		return errors.New("查询用户失败：" + err.Error())
	}
	// 检查旧密码是否正确
	if !crypto.CheckPasswordHash(in.OldPassword, user.Password) {
		return errors.New(commonModel.PASSWORD_ERROR)
	}
	if in.NewPassword == in.OldPassword {
		return errors.New("新密码不能与旧密码相同")
	}
	newPasswordHash, err := crypto.HashPassword(in.NewPassword)
	if err != nil {
		return errors.New("密码加密失败：" + err.Error())
	}
	return u.userRepo.ChangePassword(ctx, user.ID, newPasswordHash)
}

const UsernameChangeIntervalDays = 180 * 24 * time.Hour

func (u *UserService) ChangeUserId(ctx context.Context, id int64, in req.ChangeUserIdReq) error {
	user, err := u.userRepo.GetUserById(ctx, id)
	if err != nil {
		return errors.New("查询用户失败：" + err.Error())
	}
	// 检查上次更改过了多长时间
	if user.UserIdChangedAt != nil {
		gap := time.Since(*user.UserIdChangedAt)
		if gap < UsernameChangeIntervalDays {
			remainingDays := int((UsernameChangeIntervalDays - gap).Hours() / 24)
			return fmt.Errorf("短期内只能更改一次，请在 %d 天后再试", remainingDays)
		}
	}
	// 检查新用户名是否已存在
	existsUser, err := u.userRepo.GetUserByUserId(ctx, in.NewUserId)
	if err == nil && existsUser.ID != user.ID {
		return errors.New(commonModel.USERNAME_HAS_EXISTS)
	}
	now := time.Now()
	user.UserIdChangedAt = &now
	return u.userRepo.ChangeUserId(ctx, id, in.NewUserId, user.UserIdChangedAt)
}

func (u *UserService) GetUserInfo(ctx context.Context, id int64) (dto.UserInfo, error) {
	user, err := u.userRepo.GetUserById(ctx, id)
	if err != nil {
		return dto.UserInfo{}, errors.New("查询用户信息失败：" + err.Error())
	}
	return dto.UserInfo{
		UserId:    user.UserId,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Telephone: user.Telephone,
	}, nil
}

func (u *UserService) UploadAvatar(ctx context.Context, userID int64, contentType string, r io.Reader) (string, error) {
	if userID <= 0 {
		return "", errors.New("invalid userID")
	}
	if u.oss == nil {
		return "", errors.New("oss not configured")
	}
	if r == nil {
		return "", errors.New("empty file")
	}

	// 生成对象 key：avatars/{userID}/{snowflake}.bin
	// contentType 可在 handler 端按真实文件类型传入（image/jpeg、image/png）
	objKey := path.Join("avatars", fmt.Sprintf("%d", userID), snowflake.GenStringID())
	url, err := u.oss.UploadFile(ctx, objKey, strings.TrimSpace(contentType), r)
	if err != nil {
		return "", err
	}

	user, err := u.userRepo.GetUserById(ctx, userID)
	if err != nil {
		return "", errors.New("查询用户失败：" + err.Error())
	}
	user.Avatar = url
	if err := u.userRepo.UpdateUser(ctx, user); err != nil {
		return "", errors.New("更新头像失败：" + err.Error())
	}
	return url, nil
}

func NewUserService(userRepo repo.UserRepoInterface, oss ossutil.OSSInterface) *UserService {
	return &UserService{
		userRepo: userRepo,
		oss:      oss,
	}
}
