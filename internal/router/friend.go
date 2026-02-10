package router

import "minichat/internal/di"

func setupFriendRoutes(appRouterGroup *AppRouterGroup, h *di.HandlerProvider) {
	// 更改好友备注
	appRouterGroup.AuthRouterGroup.PUT("/friend/update/remark", h.FriendHandler.UpdateFriendRemark)
	// 获取好友列表
	appRouterGroup.AuthRouterGroup.GET("/friend/list", h.FriendHandler.GetFriendList)
	// 获取拉黑好友列表
	appRouterGroup.AuthRouterGroup.GET("/friend/blacklist", h.FriendHandler.GetBlackFriendList)
	// 删除好友
	appRouterGroup.AuthRouterGroup.DELETE("/friend/delete", h.FriendHandler.DeleteFriend)
	// 拉黑好友
	appRouterGroup.AuthRouterGroup.PUT("/friend/black", h.FriendHandler.BlackFriend)
	// 取消拉黑好友
	appRouterGroup.AuthRouterGroup.PUT("/friend/unblack", h.FriendHandler.UnBlackFriend)
}
