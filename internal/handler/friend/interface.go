package friend

import "github.com/gin-gonic/gin"

type FriendHandlerInterface interface {
	//拉黑好友
	BlackFriend(c *gin.Context)
	//取消拉黑好友
	UnBlackFriend(c *gin.Context)
	//删除好友
	DeleteFriend(c *gin.Context)
	//获取好友列表
	GetFriendList(c *gin.Context)
	// 修改好友备注
	UpdateFriendRemark(c *gin.Context)
	// 获取拉黑好友列表
	GetBlackFriendList(c *gin.Context)
}
