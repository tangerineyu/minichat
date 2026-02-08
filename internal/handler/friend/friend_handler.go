package friend

import (
	"minichat/internal/service/friend"
)

type FriendHandler struct {
	friendService friend.FriendServiceInterface
}

func NewFriendHandler(svc friend.FriendServiceInterface) *FriendHandler {
	return &FriendHandler{
		friendService: svc,
	}
}
