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
}
