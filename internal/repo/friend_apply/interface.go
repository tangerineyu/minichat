package friend_apply

import (
	"context"
	"minichat/internal/model"
)

type FriendApplyRepoInterface interface {
	SendFriendApply(ctx context.Context, formUserId, toUsrId, applyMsg string) (*model.FriendApply, error)
}
