package friend_apply

import (
	"context"
	"minichat/internal/model"
	"minichat/util/snowflake"

	"gorm.io/gorm"
)

var _ FriendApplyRepoInterface = (*FriendApplyRepo)(nil)

type FriendApplyRepo struct {
	db *gorm.DB
}

func (f *FriendApplyRepo) UpdateApplyStatus(ctx context.Context, applyId int64, status int) error {
	return f.db.WithContext(ctx).Model(model.FriendApply{}).
		Where("id = ?", applyId).
		Update("status", status).Error
}

func (f *FriendApplyRepo) GetFriendApplyById(ctx context.Context, applyId int64) (*model.FriendApply, error) {
	var apply model.FriendApply
	if err := f.db.WithContext(ctx).First(&apply).Where("id = ?", applyId).Error; err != nil {
		return nil, err
	}
	return &apply, nil
}

func (f *FriendApplyRepo) CreateFriendApply(ctx context.Context, formUserId, toUsrId int64, applyMsg string) (*model.FriendApply, error) {
	Id := snowflake.GenInt64ID()
	apply := &model.FriendApply{
		ID:         Id,
		FromUserId: formUserId,
		ToUserId:   toUsrId,
		ApplyMsg:   applyMsg,
		Status:     0,
	}
	if err := f.db.WithContext(ctx).Create(apply).Error; err != nil {
		return nil, err
	}
	return apply, nil
}

func NewFriendApplyRepo(db *gorm.DB) FriendApplyRepoInterface {
	return &FriendApplyRepo{db: db}
}
