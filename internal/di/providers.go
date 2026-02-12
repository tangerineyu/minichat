package di

import (
	friendHandler "minichat/internal/handler/friend"
	friendApplyHandler "minichat/internal/handler/friend_apply"
	groupHandler "minichat/internal/handler/group"
	userHandler "minichat/internal/handler/user"
)

type HandlerProvider struct {
	UserHandler        *userHandler.UserHandler
	FriendApplyHandler *friendApplyHandler.FriendApplyHandler
	FriendHandler      *friendHandler.FriendHandler
	GroupHandler       *groupHandler.GroupHandler
}
