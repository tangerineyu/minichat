package router

import "minichat/internal/di"

func SetupGroupRoutes(appRouterGroup *AppRouterGroup, h *di.HandlerProvider) {
	if h == nil || h.GroupHandler == nil {
		return
	}
	// group should be authenticated
	appRouterGroup.AuthRouterGroup.POST("/group/create", h.GroupHandler.CreateGroup)
	appRouterGroup.AuthRouterGroup.POST("/group/dissolve/:group_id", h.GroupHandler.DissolveGroup)
	appRouterGroup.AuthRouterGroup.POST("/group/update_info/:group_id", h.GroupHandler.UpdateGroupInfo)
	appRouterGroup.AuthRouterGroup.GET("/group/info", h.GroupHandler.GetGroupInfo)
	appRouterGroup.AuthRouterGroup.GET("/group/list", h.GroupHandler.GetGroupList)
	appRouterGroup.AuthRouterGroup.POST("/group/add_member", h.GroupHandler.AddGroupMember)
	appRouterGroup.AuthRouterGroup.POST("/group/delete_member", h.GroupHandler.DeleteGroupMember)
}
