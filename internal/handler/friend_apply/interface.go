package friend_apply

import "github.com/gin-gonic/gin"

type FriendApplyHandlerInterface interface {
	// 发送好友申请
	SendFriendApply(c *gin.Context)
	// 处理好友申请，接受或拒绝
	DealWithFriendApply(c *gin.Context)
}
