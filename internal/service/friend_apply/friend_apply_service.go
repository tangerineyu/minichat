package friend_apply

import (
	"context"
	"errors"
	"minichat/internal/common"
	"minichat/internal/dto"
	friend_repo "minichat/internal/repo/friend"
	repo "minichat/internal/repo/friend_apply"
	"minichat/internal/repo/user"
	"minichat/internal/req"

	"go.uber.org/zap"
)

var _ FriendApplyServiceInterface = (*FriendApplyService)(nil)

type FriendApplyService struct {
	friendApplyRepo repo.FriendApplyRepoInterface
	friendRepo      friend_repo.FriendRepoInterface
	userRepo        user.UserRepoInterface
}

func (s *FriendApplyService) GetFriendApply(ctx context.Context, id int64) ([]*dto.FriendApplyItem, error) {
	list, err := s.friendApplyRepo.GetFriendApplyList(ctx, id)
	if err != nil {
		zap.L().Error("GetFriendApply", zap.Error(err))
		return nil, errors.New("查询好友申请列表失败")
	}
	return list, nil
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
	// 如果备注为空，默认使用对方的昵称作为备注
	if in.Remark == "" {
		user1, _ := s.userRepo.GetUserById(ctx, apply.ToUserId)
		in.Remark = user1.Nickname
	}
	// 处理人对申请人的备注
	accept2RequestRemark := in.Remark
	user2, _ := s.userRepo.GetUserById(ctx, apply.FromUserId)
	// 申请人对处理人的备注，申请人不在操作界面，直接使用处理人的昵称作为备注
	request2AcceptRemark := user2.Nickname
	//事务处理
	return s.friendRepo.MakeFriends(ctx, applyId, accept2RequestRemark, request2AcceptRemark)
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
) *FriendApplyService {
	return &FriendApplyService{
		friendApplyRepo: friendApplyRepo,
		friendRepo:      friendRepo,
		userRepo:        userRepo,
	}
}
