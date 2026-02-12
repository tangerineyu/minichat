package group

import (
	"context"
	"errors"
	"minichat/internal/common"
	"minichat/internal/dto"
	"minichat/internal/model"
	groupRepo "minichat/internal/repo/group"
	groupMemberRepo "minichat/internal/repo/group_member"
	"minichat/internal/req"
	"minichat/util/snowflake"
)

var _ GroupServiceInterface = (*GroupService)(nil)

type GroupService struct {
	groupRepo       groupRepo.GroupRepoInterface
	groupMemberRepo groupMemberRepo.GroupMemberRepoInterface
}

func (g *GroupService) UpdateGroupInfo(ctx context.Context, groupId int64, operatorId int64, in req.UpdateGroupInfoReq) (*dto.GroupInfo, error) {
	// 校验操作者权限（必须在群内且非普通成员）
	operator, err := g.groupMemberRepo.GetMemberById(ctx, operatorId, groupId)
	if err != nil {
		return nil, errors.New(common.GROUP_NOT_EXISTS_USER)
	}
	if operator.Role == 1 {
		return nil, errors.New(common.GROUP_UPDATE_NO_PERMISSION)
	}

	//参数校验
	if len(in.Name) > 10 {
		return nil, errors.New(common.GROUP_NAME_TOO_LONG)
	}
	if len(in.Announcement) > 200 {
		return nil, errors.New(common.GROUP_ANNOUNCEMENT_TOO_LONG)
	}

	//获取群信息（用于校验是否已解散 + 填充默认值）
	group, err := g.groupRepo.GetGroupById(ctx, groupId)
	if err != nil {
		return nil, errors.New(common.GROUP_IS_NOT_EXIST)
	}
	if group.Status != 0 {
		return nil, errors.New(common.GROUP_IS_DISSOLVED)
	}

	//填充默认值：空值表示“不修改”
	name := in.Name
	if name == "" {
		name = group.Name
	}
	announcement := in.Announcement
	if announcement == "" {
		announcement = group.Announcement
	}
	avatar := in.Avatar
	if avatar == "" {
		avatar = group.Avatar
	}

	//更新
	if err := g.groupRepo.UpdateGroupInfo(ctx, groupId, map[string]interface{}{
		"name":         name,
		"announcement": announcement,
		"avatar":       avatar,
	}); err != nil {
		return nil, err
	}

	//返回更新后的数据：已经有最终值，不必再查一次 DB
	return &dto.GroupInfo{
		GroupId:      groupId,
		GroupName:    name,
		Announcement: announcement,
		Avatar:       avatar,
	}, nil
}

func (g *GroupService) DissolveGroup(ctx context.Context, ownerId int64, groupId int64) error {
	group, err := g.groupRepo.GetGroupById(ctx, groupId)
	if err != nil {
		return errors.New(common.GROUP_IS_NOT_EXIST)
	}
	if group.OwnerID != ownerId {
		return errors.New(common.GROUP_DISSOLVE_NO_PERMISSION)
	}
	return g.groupRepo.DissolveGroup(ctx, groupId)
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

func NewGroupService(groupRepo groupRepo.GroupRepoInterface, groupMemberRepo groupMemberRepo.GroupMemberRepoInterface) GroupServiceInterface {
	return &GroupService{
		groupRepo:       groupRepo,
		groupMemberRepo: groupMemberRepo,
	}
}
