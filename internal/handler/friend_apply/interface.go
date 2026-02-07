package friend_apply

import "github.com/gin-gonic/gin"

type FriendApplyHandlerInterface interface {
	SendFriendApply(c *gin.Context)
}
