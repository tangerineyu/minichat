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
