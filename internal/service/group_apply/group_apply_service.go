package group_apply

import (
	"context"
	"errors"
	"minichat/internal/common"
	"minichat/internal/model"
	groupApplyRepo "minichat/internal/repo/group_apply"
	groupMemberRepo "minichat/internal/repo/group_member"
)

var _ GroupApplyServiceInterface = (*GroupApplyService)(nil)

type GroupApplyService struct {
	groupApplyRepo  groupApplyRepo.GroupApplyRepoInterface
	groupMemberRepo groupMemberRepo.GroupMemberRepoInterface
}

func NewGroupApplyService(groupApplyRepo groupApplyRepo.GroupApplyRepoInterface, groupMemberRepo groupMemberRepo.GroupMemberRepoInterface) GroupApplyServiceInterface {
	return &GroupApplyService{groupApplyRepo: groupApplyRepo, groupMemberRepo: groupMemberRepo}
}

func (g *GroupApplyService) validateOperator(ctx context.Context, operatorId int64, groupId int64) error {
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

func (g *GroupApplyService) DealGroupApply(ctx context.Context, operatorId int64, groupId int64, applyId int64, accept int8) error {
	// 校验操作者权限（必须在群内且非普通成员）
	if err := g.validateOperator(ctx, operatorId, groupId); err != nil {
		return err
	}

	// 获取申请信息
	apply, err := g.groupApplyRepo.GetGroupApplyById(ctx, applyId)
	if err != nil || apply == nil || apply.GroupID != groupId {
		return errors.New(common.GROUP_APPLY_IS_NOT_EXIST)
	}
	if apply.Status != 0 {
		return errors.New(common.GROUP_APPLY_ALREADY_PROCESSED)
	}

	// 更新申请状态
	err = g.groupApplyRepo.HandleGroupApply(ctx, applyId, operatorId, accept)
	if err != nil {
		return err
	}
	return nil
}

func (g *GroupApplyService) GetGroupApplyList(ctx context.Context, userId int64) ([]*model.GroupApply, error) {
	return g.groupApplyRepo.GetGroupApplyList(ctx, userId)
}

func (g *GroupApplyService) GetGroupApplyListByGroupId(ctx context.Context, operatorId int64, groupId int64) ([]*model.GroupApply, error) {
	if err := g.validateOperator(ctx, operatorId, groupId); err != nil {
		return nil, err
	}
	return g.groupApplyRepo.GetGroupApplyListByGroupId(ctx, groupId)
}

func (g *GroupApplyService) GetGroupApplyDetail(ctx context.Context, operatorId int64, applyId int64) (*model.GroupApply, error) {
	apply, err := g.groupApplyRepo.GetGroupApplyById(ctx, applyId)
	if err != nil || apply == nil {
		return nil, errors.New(common.GROUP_APPLY_IS_NOT_EXIST)
	}

	// 申请本人 / 邀请人：可以看
	if apply.UserID == operatorId || apply.InviterId == operatorId {
		return apply, nil
	}

	// 管理员/群主：可以看
	if err := g.validateOperator(ctx, operatorId, apply.GroupID); err != nil {
		return nil, err
	}

	return apply, nil
}
