package group

import (
	"context"
	"minichat/internal/dto"
	"minichat/internal/req"
)

type GroupServiceInterface interface {
	CreateGroup(ctx context.Context, ownerId int64, in req.CreateGroupReq) (*dto.GroupInfo, error)

	DissolveGroup(ctx context.Context, ownerId int64, groupId int64) error

	UpdateGroupInfo(ctx context.Context, groupId int64, operatorId int64, in req.UpdateGroupInfoReq) (*dto.GroupInfo, error)
	// 获取群信息，包含成员列表和成员角色
	GetGroupInfo(ctx context.Context, groupId int64, userId int64) (*dto.GroupDetailInfo, error)
	// 添加群成员，返回添加结果
	AddGroupMembers(ctx context.Context, operatorId int64, groupId int64, userIds []int64) (*dto.GroupAddStatus, error)
}
