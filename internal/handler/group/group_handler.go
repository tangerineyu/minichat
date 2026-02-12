package group

import (
	"minichat/internal/handler/response"
	"minichat/internal/req"
	service "minichat/internal/service/group"

	"github.com/gin-gonic/gin"
)

var _ GroupHandlerInterface = (*GroupHandler)(nil)

type GroupHandler struct {
	groupService service.GroupServiceInterface
}

func (g GroupHandler) CreateGroup(c *gin.Context) {
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

func (g GroupHandler) DissolveGroup(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (g GroupHandler) UpdateGroupInfo(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (g GroupHandler) GetGroupInfo(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (g GroupHandler) GetGroupList(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (g GroupHandler) AddGroupMember(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (g GroupHandler) DeleteGroupMember(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (g GroupHandler) GetGroupMemberList(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func NewGroupHandler(groupService service.GroupServiceInterface) *GroupHandler {
	return &GroupHandler{groupService: groupService}
}
