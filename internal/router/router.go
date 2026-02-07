package router

import (
	"minichat/internal/di"
	"minichat/internal/middlerware"

	"github.com/gin-gonic/gin"
)

type AppRouterGroup struct {
	PublicRouterGroup *gin.RouterGroup
	AuthRouterGroup   *gin.RouterGroup
	WSRouterGroup     *gin.RouterGroup
}

// SetupRouter 配置路由
func SetupRouter(r *gin.Engine, h *di.HandlerProvider) {
	setupMiddleware(r)
	appRouterGroup := setupRouterGroup(r)
	setupUserRoutes(appRouterGroup, h)

	setupFriendApplyRoutes(appRouterGroup, h)
}

func setupRouterGroup(r *gin.Engine) *AppRouterGroup {
	publicRouterGroup := r.Group("/api")
	authRouterGroup := r.Group("/api")
	authRouterGroup.Use(middlerware.JWTAuthMiddleware())
	wsRouterGroup := r.Group("/ws")
	return &AppRouterGroup{
		PublicRouterGroup: publicRouterGroup,
		AuthRouterGroup:   authRouterGroup,
		WSRouterGroup:     wsRouterGroup,
	}
}
