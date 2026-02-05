package router

import (
	"minichat/internal/di"
)

// RegisterUserRoutes registers user-related endpoints.
func setupUserRoutes(appRouterGroup *AppRouterGroup, h *di.HandlerProvider) {
	appRouterGroup.PublicRouterGroup.POST("/login", h.UserHandler.Login)
	appRouterGroup.PublicRouterGroup.POST("/register", h.UserHandler.Register)
	appRouterGroup.AuthRouterGroup.PUT("/user/update/info", h.UserHandler.UpdateUserInfo)
	appRouterGroup.AuthRouterGroup.PUT("/user/update/password", h.UserHandler.ChangePassword)
	appRouterGroup.AuthRouterGroup.PUT("/user/update/username", h.UserHandler.ChangeUsername)
}
