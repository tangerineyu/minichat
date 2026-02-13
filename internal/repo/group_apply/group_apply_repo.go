package group_apply

import (
	"context"
	"minichat/internal/model"

	"gorm.io/gorm"
)

type GroupApplyRepo struct {
	db *gorm.DB
}

func (g *GroupApplyRepo) CreateGroupApply(ctx context.Context, formUserId, groupId int64, reason string) (*model.GroupApply, error) {
	var groupApply model.GroupApply
	err := g.db.WithContext(ctx).Model(&groupApply).Create(
		&model.GroupApply{
			UserID:  formUserId,
			GroupID: groupId,
			Reason:  reason,
			Status:  0,
		},
	).Error
	if err != nil {
		return nil, err
	}
	return &groupApply, nil
}

func (g *GroupApplyRepo) HandleGroupApply(ctx context.Context, applyId int64, operatorId int64, status int) error {
	//TODO implement me
	panic("implement me")
}

func (g *GroupApplyRepo) GetGroupApplyList(ctx context.Context, userId int64) ([]*model.GroupApply, error) {
	//TODO implement me
	panic("implement me")
}

func NewGroupApplyRepo(db *gorm.DB) GroupApplyRepoInterface {
	return &GroupApplyRepo{db: db}
}
