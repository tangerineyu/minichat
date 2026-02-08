package friend_apply

import (
	"context"
	"errors"
	"minichat/internal/common"
	friend_repo "minichat/internal/repo/friend"
	repo "minichat/internal/repo/friend_apply"
	"minichat/internal/repo/user"
	"minichat/internal/req"
)

var _ FriendApplyServiceInterface = (*FriendApplyService)(nil)

type FriendApplyService struct {
	friendApplyRepo repo.FriendApplyRepoInterface
	friendRepo      friend_repo.FriendRepoInterface
	userRepo        user.UserRepoInterface
}

func (s *FriendApplyService) DealWithFriendApply(ctx context.Context, Id int64, in req.DealWithFriendApplyReq) error {
	applyId := in.ApplyId

	apply, err := s.friendApplyRepo.GetFriendApplyById(ctx, applyId)
	if err != nil {
		return errors.New(common.APPLY_NOT_FOUND)
	}
	//这里userId (string)是类似微信号的用户可自定义的字符串，而apply.ToUserId（int64）是数据库中存储的用户ID，
	//需要根据apply.ToUserId查询用户信息，获取到用户的userId，再进行比较
	//申请不属于该用户
	if Id != apply.ToUserId {
		return errors.New(common.APPLY_NOT_BELONG_TO_USER)
	}
	//申请已经处理
	if apply.Status != 0 {
		return errors.New(common.APPLY_ALREADY_DEAL)
	}
	action := in.Action
	if action == 2 {
		return s.friendApplyRepo.UpdateApplyStatus(ctx, applyId, 2)
	}
	//TODO：同意好友申请后，添加好友关系
	//事务处理
	return s.friendRepo.MakeFriends(ctx, applyId)
}

func (s *FriendApplyService) SendFriendApply(ctx context.Context, fromUserId int64, in req.SendFriendApplyReq) error {
	if fromUserId == in.ToUserId {
		return errors.New(common.APPLY_FRIEND_SELF)
	}
	_, err := s.friendApplyRepo.CreateFriendApply(ctx, fromUserId, in.ToUserId, in.Message)
	return err
}

func NewFriendApplyService(
	friendApplyRepo repo.FriendApplyRepoInterface,
	friendRepo friend_repo.FriendRepoInterface,
	userRepo user.UserRepoInterface,
) FriendApplyServiceInterface {
	return &FriendApplyService{
		friendApplyRepo: friendApplyRepo,
		friendRepo:      friendRepo,
		userRepo:        userRepo,
	}
}
