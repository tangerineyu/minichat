package group_apply

import (
	"context"
	"minichat/internal/model"
)

type GroupApplyServiceInterface interface {
	// 处理群成员申请，接受或拒绝
	DealGroupApply(ctx context.Context, operatorId int64, groupId int64, applyId int64, accept int8) error
	// 获取用户相关的加群申请列表（发起者是自己 or 邀请者是自己）
	GetGroupApplyList(ctx context.Context, userId int64) ([]*model.GroupApply, error)
	// 获取群的加群申请列表（仅管理员/群主可看）
	GetGroupApplyListByGroupId(ctx context.Context, operatorId int64, groupId int64) ([]*model.GroupApply, error)
	// 获取加群申请详情（本人相关或管理员/群主可看）
	GetGroupApplyDetail(ctx context.Context, operatorId int64, applyId int64) (*model.GroupApply, error)
}
