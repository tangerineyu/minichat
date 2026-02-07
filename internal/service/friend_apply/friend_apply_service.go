package friend_apply

import (
	"context"
	"errors"
	"minichat/internal/common"
	repo "minichat/internal/repo/friend_apply"
	"minichat/internal/req"
)

type FriendApplyService struct {
	friendApplyRepo repo.FriendApplyRepoInterface
}

var _ FriendApplyServiceInterface = (*FriendApplyService)(nil)

func (s *FriendApplyService) SendFriendApply(ctx context.Context, fromUserId string, in req.SendFriendApplyReq) error {
	if fromUserId == in.ToUserId {
		return errors.New(common.APPLY_FRIEND_SELF)
	}
	_, err := s.friendApplyRepo.SendFriendApply(ctx, fromUserId, in.ToUserId, in.Message)
	return err
}

func NewFriendApplyService(friendApplyRepo repo.FriendApplyRepoInterface) FriendApplyServiceInterface {
	return &FriendApplyService{
		friendApplyRepo: friendApplyRepo,
	}
}
