package di

import (
	friendApplyHandler "minichat/internal/handler/friend_apply"
	userHandler "minichat/internal/handler/user"
)

type HandlerProvider struct {
	UserHandler        *userHandler.UserHandler
	FriendApplyHandler *friendApplyHandler.FriendApplyHandler
}
