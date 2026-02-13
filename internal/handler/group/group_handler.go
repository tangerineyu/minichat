package group

import (
	"minichat/internal/handler/response"
	"minichat/internal/req"
	service "minichat/internal/service/group"
	"strconv"

	"github.com/gin-gonic/gin"
)

var _ GroupHandlerInterface = (*GroupHandler)(nil)

type GroupHandler struct {
	groupService service.GroupServiceInterface
}

func (g *GroupHandler) DealGroupApply(c *gin.Context) {

}

func (g *GroupHandler) CreateGroup(c *gin.Context) {
	var in req.CreateGroupReq
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, 400, "invalid request")
		return
	}
	id := c.GetInt64("id")
	groupInfo, err := g.groupService.CreateGroup(c.Request.Context(), id, in)
	if err != nil {
		response.ServerError(c, err)
		return
	}
	response.Success(c, gin.H{"group_info": groupInfo})
}

func (g *GroupHandler) DissolveGroup(c *gin.Context) {
	id := c.GetInt64("id")
	groupIdStr := c.Param("group_id")
	groupId, _ := strconv.ParseInt(groupIdStr, 10, 64)
	err := g.groupService.DissolveGroup(c.Request.Context(), id, groupId)
	if err != nil {
		response.Fail(c, 403, err.Error())
		return
	}
	response.Success(c, "群组已解散")
}

func (g *GroupHandler) UpdateGroupInfo(c *gin.Context) {
	id := c.GetInt64("id")
	groupIdStr := c.Param("group_id")
	groupId, _ := strconv.ParseInt(groupIdStr, 10, 64)

	var in req.UpdateGroupInfoReq
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, 400, "invalid request")
		return
	}
	updatedGroupInfo, err := g.groupService.UpdateGroupInfo(c.Request.Context(), groupId, id, in)
	if err != nil {
		response.Fail(c, 403, err.Error())
		return
	}
	response.Success(c, updatedGroupInfo)
}

func (g *GroupHandler) GetGroupInfo(c *gin.Context) {
	groupIdStr := c.Param("group_id")
	groupId, _ := strconv.ParseInt(groupIdStr, 10, 64)
	groupInfo, err := g.groupService.GetGroupInfo(c.Request.Context(), groupId, c.GetInt64("id"))
	if err != nil {
		response.Fail(c, 403, err.Error())
		return
	}
	response.Success(c, groupInfo)
}

func (g *GroupHandler) GetGroupList(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupHandler) AddGroupMembers(c *gin.Context) {
	id := c.GetInt64("id")
	groupIdStr := c.Param("group_id")
	groupId, _ := strconv.ParseInt(groupIdStr, 10, 64)
	var in req.AddGroupMembersReq
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, 400, "invalid request")
		return
	}
	status, err := g.groupService.AddGroupMembers(c.Request.Context(), id, groupId, in.UserIds)
	if err != nil {
		response.Fail(c, 403, err.Error())
		return
	}
	response.Success(c, status)
}

func (g *GroupHandler) DeleteGroupMember(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupHandler) GetGroupMemberList(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func NewGroupHandler(groupService service.GroupServiceInterface) *GroupHandler {
	return &GroupHandler{groupService: groupService}
}
