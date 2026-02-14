package group_apply

import (
	"context"
	"minichat/internal/model"

	"gorm.io/gorm"
)

type GroupApplyRepo struct {
	db *gorm.DB
}

func (g *GroupApplyRepo) GetGroupApplyListByGroupId(ctx context.Context, groupId int64) ([]*model.GroupApply, error) {
	var applies []*model.GroupApply
	err := g.db.WithContext(ctx).
		Model(&model.GroupApply{}).
		Where("group_id = ?", groupId).
		Order("created_at desc").
		Find(&applies).Error
	if err != nil {
		return nil, err
	}
	return applies, nil
}

func (g *GroupApplyRepo) GetGroupApplyById(ctx context.Context, applyId int64) (*model.GroupApply, error) {
	var groupApply *model.GroupApply
	err := g.db.WithContext(ctx).Where("id = ?", applyId).First(&groupApply).Error
	return groupApply, err
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

// 处理群申请，status: 1同意，2拒绝
func (g *GroupApplyRepo) HandleGroupApply(ctx context.Context, applyId int64, operatorId int64, status int8) error {
	return g.db.WithContext(ctx).Model(model.GroupApply{}).
		Where("id = ?", applyId).
		Update("operator_id", operatorId).
		Update("status", status).Error
}

func (g *GroupApplyRepo) GetGroupApplyList(ctx context.Context, userId int64) ([]*model.GroupApply, error) {
	var applies []*model.GroupApply
	err := g.db.WithContext(ctx).
		Model(&model.GroupApply{}).
		Where("user_id = ? OR inviter_id = ?", userId, userId).
		Order("created_at desc").
		Find(&applies).Error
	if err != nil {
		return nil, err
	}
	return applies, nil
}

func NewGroupApplyRepo(db *gorm.DB) GroupApplyRepoInterface {
	return &GroupApplyRepo{db: db}
}
