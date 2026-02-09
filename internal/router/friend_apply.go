package router

import "minichat/internal/di"

func setupFriendApplyRoutes(appRouterGroup *AppRouterGroup, h *di.HandlerProvider) {
	if h == nil || h.FriendApplyHandler == nil {
		return
	}
	// friend apply should be authenticated
	appRouterGroup.AuthRouterGroup.POST("/friend/apply/send", h.FriendApplyHandler.SendFriendApply)
	appRouterGroup.AuthRouterGroup.POST("/friend/apply/deal", h.FriendApplyHandler.DealWithFriendApply)
	appRouterGroup.AuthRouterGroup.GET("/friend/apply/list", h.FriendApplyHandler.GetFriendApplyList)
}
