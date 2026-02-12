package group

import (
	"context"
	"minichat/internal/model"
)

type GroupRepoInterface interface {
	//创建群
	CreateGroup(ctx context.Context, info *model.Group) error
	// 根据ID获取群信息
	GetGroupById(ctx context.Context, groupId int64) (*model.Group, error)
	// 解散群
	DissolveGroup(ctx context.Context, groupId int64) error
	// 更新群信息
	UpdateGroupInfo(ctx context.Context, groupId int64, in map[string]interface{}) error
}
