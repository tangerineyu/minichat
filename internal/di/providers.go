package di

import (
	friendHandler "minichat/internal/handler/friend"
	friendApplyHandler "minichat/internal/handler/friend_apply"
	groupHandler "minichat/internal/handler/group"
	groupApplyHandler "minichat/internal/handler/group_apply"
	userHandler "minichat/internal/handler/user"
)

type HandlerProvider struct {
	UserHandler        *userHandler.UserHandler
	FriendApplyHandler *friendApplyHandler.FriendApplyHandler
	FriendHandler      *friendHandler.FriendHandler
	GroupHandler       *groupHandler.GroupHandler
	GroupApplyHandler  *groupApplyHandler.GroupApplyHandler
}
