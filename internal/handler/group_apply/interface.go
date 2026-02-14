package group_apply

import "github.com/gin-gonic/gin"

type GroupApplyHandlerInterface interface {
	//处理群成员申请
	DealGroupApply(c *gin.Context)
	//获取用户的加群申请列表
	GetUserGroupApplyList(c *gin.Context)
	//获取群的加群申请列表
	GetGroupApplyListByGroupId(c *gin.Context)
	//获取加群申请详情
	GetGroupApplyDetail(c *gin.Context)
}
