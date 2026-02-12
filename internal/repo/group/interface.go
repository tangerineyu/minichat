package group

import (
	"context"
	"minichat/internal/model"
)

type GroupRepoInterface interface {
	//创建群
	CreateGroup(ctx context.Context, info *model.Group) error
}
