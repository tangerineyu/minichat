package group

import (
	"context"
	"minichat/internal/model"

	"gorm.io/gorm"
)

var _ GroupRepoInterface = (*GroupRepo)(nil)

type GroupRepo struct {
	db *gorm.DB
}

func (g *GroupRepo) UpdateGroupInfo(ctx context.Context, groupId int64, in map[string]interface{}) error {
	return g.db.WithContext(ctx).Model(&model.Group{}).Where("id=?", groupId).Updates(in).Error
}

func (g *GroupRepo) GetGroupById(ctx context.Context, groupId int64) (*model.Group, error) {
	var group *model.Group
	if err := g.db.WithContext(ctx).Where("id = ?", groupId).First(&group).Error; err != nil {
		return nil, err
	}
	return group, nil
}

func (g *GroupRepo) DissolveGroup(ctx context.Context, groupId int64) error {
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Group{}).Where("id = ?", groupId).Update("status", 1).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.GroupMember{}).Where("group_id = ?", groupId).Update("status", 2).Error; err != nil {
			return err
		}
		return nil
	})
}

func (g *GroupRepo) CreateGroup(ctx context.Context, info *model.Group) error {
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(info).Error; err != nil {
			return err
		}
		member := &model.GroupMember{
			GroupID: info.ID,
			UserID:  info.OwnerID,
			Role:    2, // Owner
		}
		return tx.Create(member).Error
	})
}

func NewGroupRepo(db *gorm.DB) *GroupRepo {
	return &GroupRepo{db: db}
}
