package friend_apply

import (
	"context"
	"minichat/internal/model"

	"gorm.io/gorm"
)

var _ FriendApplyRepoInterface = (*FriendApplyRepo)(nil)

type FriendApplyRepo struct {
	db *gorm.DB
}

func (f FriendApplyRepo) SendFriendApply(ctx context.Context, formUserId, toUsrId, applyMsg string) (*model.FriendApply, error) {
	apply := &model.FriendApply{
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
