package group_member

import (
	"context"
	"minichat/internal/model"
)

type GroupMemberRepoInterface interface {
	GetMemberById(ctx context.Context, memberId int64, groupId int64) (*model.GroupMember, error)
}
