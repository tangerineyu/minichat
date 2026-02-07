package friend_apply

import (
	"context"
	"minichat/internal/req"
)

type FriendApplyServiceInterface interface {
	SendFriendApply(ctx context.Context, fromUserId string, in req.SendFriendApplyReq) error
}
