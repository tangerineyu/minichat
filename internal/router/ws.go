package router

import "minichat/internal/di"

func setupWSRoutes(appRouterGroup *AppRouterGroup, h *di.HandlerProvider) {
	if appRouterGroup == nil || appRouterGroup.WSRouterGroup == nil || h == nil || h.WSHandler == nil {
		return
	}

	// WebSocket chat入口：
	//   - Header: Authorization: Bearer <access_token>
	//   - 或 query: /ws/chat?token=<access_token>
	appRouterGroup.WSRouterGroup.GET("/chat", h.WSHandler.Connect)
}
