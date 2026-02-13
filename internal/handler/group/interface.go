package group

import "github.com/gin-gonic/gin"

type GroupHandlerInterface interface {
	//创建群
	CreateGroup(c *gin.Context)
	//解散群
	DissolveGroup(c *gin.Context)
	//修改群信息
	UpdateGroupInfo(c *gin.Context)
	//获取群信息
	GetGroupInfo(c *gin.Context)
	//获取群列表
	GetGroupList(c *gin.Context)
	//添加群成员
	AddGroupMembers(c *gin.Context)
	//删除群成员
	DeleteGroupMember(c *gin.Context)
	//获取群成员列表
	GetGroupMemberList(c *gin.Context)
	//处理群成员申请
	DealGroupApply(c *gin.Context)
}
