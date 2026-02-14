package router

import "minichat/internal/di"

func setupGroupApplyRoutes(appRouterGroup *AppRouterGroup, h *di.HandlerProvider) {
	if h == nil || h.GroupApplyHandler == nil {
		return
	}
	// group apply should be authenticated
	appRouterGroup.AuthRouterGroup.POST("/group/apply/deal/:group_id", h.GroupApplyHandler.DealGroupApply)
	appRouterGroup.AuthRouterGroup.GET("/group/apply/userlist", h.GroupApplyHandler.GetUserGroupApplyList)
	appRouterGroup.AuthRouterGroup.GET("/group/apply/group_list/:group_id", h.GroupApplyHandler.GetGroupApplyListByGroupId)
	appRouterGroup.AuthRouterGroup.GET("/group/apply/group_detail/:apply_id", h.GroupApplyHandler.GetGroupApplyDetail)
}
