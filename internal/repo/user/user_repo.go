package user

import (
	"context"
	"minichat/internal/model"
	"time"

	"gorm.io/gorm"
)

// 确保UserRepo实现了UserRepoInterface接口
var _ UserRepoInterface = (*UserRepo)(nil)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepo) GetUserById(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return &user, err
	}
	return &user, nil
}

func (r *UserRepo) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetUserByTelephone(ctx context.Context, telephone string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("telephone = ?", telephone).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *UserRepo) CheckUserExistsByTelephone(ctx context.Context, telephone string) (bool, error) {
	if _, err := r.GetUserByTelephone(ctx, telephone); err != nil {
		return false, nil
	}
	return true, nil
}

func (r *UserRepo) UpdateUser(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Model(user).
		Select("avatar,nickname, updated_at").
		Updates(user).Error
}

// ChangePassword 修改用户密码
func (r *UserRepo) ChangePassword(ctx context.Context, userID int64, newHashedPassword string) error {
	return r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", userID).
		Update("password", newHashedPassword).Error
}

func (r *UserRepo) ChangeUsername(ctx context.Context, username string, usernameChangedAt *time.Time) error {
	return r.db.WithContext(ctx).Model(&model.User{}).
		Where("username = ?", username).
		Updates(map[string]interface{}{
			"username":            username,
			"username_changed_at": usernameChangedAt,
		}).Error
}
