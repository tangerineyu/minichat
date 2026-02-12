package group

import (
	"context"
	"minichat/internal/dto"
	"minichat/internal/model"
	repo "minichat/internal/repo/group"
	"minichat/internal/req"
	"minichat/util/snowflake"
)

var _ GroupServiceInterface = (*GroupService)(nil)

type GroupService struct {
	groupRepo repo.GroupRepoInterface
}

func (g *GroupService) CreateGroup(ctx context.Context, ownerId int64, in req.CreateGroupReq) (*dto.GroupInfo, error) {
	groupId := snowflake.GenInt64ID()
	newGroup := &model.Group{
		ID:           groupId,
		OwnerID:      ownerId,
		Name:         in.Name,
		Announcement: in.Announcement,
		Avatar:       in.Avatar,
	}
	err := g.groupRepo.CreateGroup(ctx, newGroup)
	if err != nil {
		return nil, err
	}

	//返回群信息
	return &dto.GroupInfo{GroupId: groupId, GroupName: in.Name, Announcement: in.Announcement, Avatar: in.Avatar}, nil
}

func NewGroupService(groupRepo repo.GroupRepoInterface) *GroupService {
	return &GroupService{
		groupRepo: groupRepo,
	}
}
