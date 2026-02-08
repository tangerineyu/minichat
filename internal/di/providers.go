package di

import (
	friendHandler "minichat/internal/handler/friend"
	friendApplyHandler "minichat/internal/handler/friend_apply"
	userHandler "minichat/internal/handler/user"
)

type HandlerProvider struct {
	UserHandler        *userHandler.UserHandler
	FriendApplyHandler *friendApplyHandler.FriendApplyHandler
	FriendHandler      *friendHandler.FriendHandler
}
