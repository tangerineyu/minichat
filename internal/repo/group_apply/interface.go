package group_apply

import (
	"context"
	"minichat/internal/model"
)

type GroupApplyRepoInterface interface {
	// 创建加入群申请
	CreateGroupApply(ctx context.Context, formUserId, groupId int64, reason string) (*model.GroupApply, error)
	// 处理加入群申请
	HandleGroupApply(ctx context.Context, applyId int64, operatorId int64, status int8) error
	// 获取用户的加入群申请列表
	GetGroupApplyList(ctx context.Context, userId int64) ([]*model.GroupApply, error)
	// 获取群的加入群申请列表
	GetGroupApplyListByGroupId(ctx context.Context, groupId int64) ([]*model.GroupApply, error)
	// 获取加入群申请详情
	GetGroupApplyById(ctx context.Context, applyId int64) (*model.GroupApply, error)
}
