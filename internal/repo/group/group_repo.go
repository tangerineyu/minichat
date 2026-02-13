package group

import (
	"context"
	"minichat/internal/model"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ GroupRepoInterface = (*GroupRepo)(nil)

type GroupRepo struct {
	db *gorm.DB
}

func (g *GroupRepo) AddGroupMembers(ctx context.Context, operatorId, groupId int64, userIds []int64) error {
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var members []model.GroupMember
		for _, uid := range userIds {
			members = append(members, model.GroupMember{
				GroupID:   groupId,
				UserID:    uid,
				InviterID: operatorId,
				Role:      0,
				Status:    0,
			})
		}
		return tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "group_id"}, {Name: "user_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"status":     0,
				"role":       0,
				"updated_at": time.Now(),
			}),
		}).Create(&members).Error
	})
}

func (g *GroupRepo) UpdateGroupInfo(ctx context.Context, groupId int64, in map[string]interface{}) error {
	return g.db.WithContext(ctx).Model(&model.Group{}).Where("id=?", groupId).Updates(in).Error
}

func (g *GroupRepo) GetGroupById(ctx context.Context, groupId int64) (*model.Group, error) {
	var group *model.Group
	if err := g.db.WithContext(ctx).Where("id = ? AND status = ?", groupId, 0).First(&group).Error; err != nil {
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
