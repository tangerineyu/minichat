package di

import (
	friendHandler "minichat/internal/handler/friend"
	friendApplyHandler "minichat/internal/handler/friend_apply"
	groupHandler "minichat/internal/handler/group"
	groupApplyHandler "minichat/internal/handler/group_apply"
	messageHandler "minichat/internal/handler/message"
	userHandler "minichat/internal/handler/user"
	wsHandler "minichat/internal/handler/ws"
)

type HandlerProvider struct {
	UserHandler        *userHandler.UserHandler
	FriendApplyHandler *friendApplyHandler.FriendApplyHandler
	FriendHandler      *friendHandler.FriendHandler
	GroupHandler       *groupHandler.GroupHandler
	GroupApplyHandler  groupApplyHandler.GroupApplyHandlerInterface
	MessageHandler     messageHandler.MessageHandlerInterface
	WSHandler          wsHandler.WSHandlerInterface
}
