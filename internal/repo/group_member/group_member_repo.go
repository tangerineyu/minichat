package group_member

import (
	"context"
	"minichat/internal/model"

	"gorm.io/gorm"
)

type GroupMemberRepo struct {
	db *gorm.DB
}

func (g *GroupMemberRepo) GetMemberById(ctx context.Context, memberId, groupId int64) (*model.GroupMember, error) {
	var groupMemberInfo model.GroupMember
	err := g.db.WithContext(ctx).Model(&model.GroupMember{}).
		Where("group_id = ? AND user_id = ?", groupId, memberId).
		First(&groupMemberInfo).Error
	if err != nil {
		return nil, err
	}
	return &groupMemberInfo, nil
}

func NewGroupMemberRepo(db *gorm.DB) GroupMemberRepoInterface {
	return &GroupMemberRepo{db: db}
}
