package router

import "minichat/internal/di"

func setupMessageRoutes(appRouterGroup *AppRouterGroup, h *di.HandlerProvider) {
	if h == nil || h.MessageHandler == nil {
		return
	}
	// message should be authenticated
	appRouterGroup.AuthRouterGroup.POST("/message/send", h.MessageHandler.SendMessage)
	appRouterGroup.AuthRouterGroup.GET("/message/history", h.MessageHandler.GetMessageHistory)
}
