package group

import (
	"context"
	"minichat/internal/dto"
	"minichat/internal/req"
)

type GroupServiceInterface interface {
	CreateGroup(ctx context.Context, ownerId int64, in req.CreateGroupReq) (*dto.GroupInfo, error)
}
