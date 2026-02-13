package group

import (
	"context"
	"errors"
	"fmt"
	"minichat/internal/common"
	"minichat/internal/dto"
	"minichat/internal/model"
	groupRepo "minichat/internal/repo/group"
	groupApplyRepo "minichat/internal/repo/group_apply"
	groupMemberRepo "minichat/internal/repo/group_member"
	"minichat/internal/req"
	"minichat/util/snowflake"

	"go.uber.org/zap"
)

var _ GroupServiceInterface = (*GroupService)(nil)

type GroupService struct {
	groupRepo       groupRepo.GroupRepoInterface
	groupMemberRepo groupMemberRepo.GroupMemberRepoInterface
	groupApplyRepo  groupApplyRepo.GroupApplyRepoInterface
}

func (g *GroupService) validateOperator(ctx context.Context, operatorId int64, groupId int64) error {
	operator, err := g.groupMemberRepo.GetMemberById(ctx, operatorId, groupId)
	if err != nil {
		return errors.New(common.GROUP_NOT_EXISTS_USER)
	}
	if operator.Status == 2 {
		return errors.New(common.GROUP_NOT_EXISTS_USER)
	}
	// 如果是管理员和群主都可以操作
	if operator.Role == 1 || operator.Role == 2 {
		return nil
	}
	return errors.New(common.GROUP_ADD_MEMBER_NO_PERMISSION)
}

func (g *GroupService) AddGroupMembers(ctx context.Context, operatorId int64, groupId int64, userIds []int64) (*dto.GroupAddStatus, error) {
	// 先确认操作者是群内成员（且未退群）
	operator, err := g.groupMemberRepo.GetMemberById(ctx, operatorId, groupId)
	if err != nil || operator == nil || operator.Status == 2 {
		return nil, errors.New(common.GROUP_NOT_EXISTS_USER)
	}

	// 管理员/群主：直接拉人进群
	if operator.Role == 1 || operator.Role == 2 {
		if err := g.groupRepo.AddGroupMembers(ctx, operatorId, groupId, userIds); err != nil {
			return nil, err
		}
		return &dto.GroupAddStatus{
			Status:       dto.GroupAddJoined,
			AddedUserIds: userIds,
			Msg:          "已直接加入群聊",
		}, nil
	}

	// 普通成员：创建入群申请，等待管理员/群主审核
	// 当前实现按 userId 逐条创建申请（不改表结构）。
	applyIds := make([]int64, 0, len(userIds))
	for _, userId := range userIds {
		reason := fmt.Sprintf("邀请用户 %d 加入群聊", userId)
		apply, err := g.groupApplyRepo.CreateGroupApply(ctx, operatorId, groupId, reason)
		if err != nil {
			zap.L().Error("创建加群申请失败", zap.Int64("operatorId", operatorId), zap.Int64("groupId", groupId), zap.Int64("userId", userId), zap.Error(err))
			return nil, fmt.Errorf("创建加群申请失败: %w", err)
		}
		if apply != nil {
			applyIds = append(applyIds, apply.ID)
		}
	}

	return &dto.GroupAddStatus{
		Status:   dto.GroupAddPendingApproval,
		ApplyIds: applyIds,
		Msg:      "已提交申请，等待管理员审核",
	}, nil
}

func (g *GroupService) GetGroupInfo(ctx context.Context, groupId int64, userId int64) (*dto.GroupDetailInfo, error) {
	group, err := g.groupRepo.GetGroupById(ctx, groupId)
	if err != nil {
		return nil, errors.New(common.GROUP_IS_NOT_EXIST)
	}
	member, _ := g.groupMemberRepo.GetMemberById(ctx, userId, groupId)
	isMember := member != nil && member.Status != 2
	groupInfo := &dto.GroupDetailInfo{
		GroupInfo: dto.GroupInfo{
			GroupId:      group.ID,
			GroupName:    group.Name,
			Announcement: group.Announcement,
			Avatar:       group.Avatar,
		},
	}
	if isMember {
		members, err := g.groupMemberRepo.GetGroupMembers(ctx, groupId)
		if err != nil {
			return nil, errors.New("获取群成员失败")
		}
		membersInfo := make([]*dto.GroupMemberInfo, 0, len(members))
		for _, m := range members {
			membersInfo = append(membersInfo, &dto.GroupMemberInfo{
				MemberId:  m.UserID,
				Nickname:  m.Nickname,
				Role:      m.Role,
				EntryTime: m.CreatedAt,
			})
		}
		groupInfo.Members = membersInfo
	}
	return groupInfo, nil
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

func NewGroupService(
	groupRepo groupRepo.GroupRepoInterface,
	groupMemberRepo groupMemberRepo.GroupMemberRepoInterface,
	groupApplyRepo groupApplyRepo.GroupApplyRepoInterface,
) GroupServiceInterface {
	return &GroupService{
		groupRepo:       groupRepo,
		groupMemberRepo: groupMemberRepo,
		groupApplyRepo:  groupApplyRepo,
	}
}
