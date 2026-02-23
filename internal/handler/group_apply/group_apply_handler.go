package group_apply

import (
	"minichat/internal/handler/response"
	"minichat/internal/req"
	groupApplyService "minichat/internal/service/group_apply"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GroupApplyHandler struct {
	groupApplyService groupApplyService.GroupApplyServiceInterface
}

func (g *GroupApplyHandler) GetUserGroupApplyList(c *gin.Context) {
	userId := c.GetInt64("id")
	list, err := g.groupApplyService.GetGroupApplyList(c.Request.Context(), userId)
	if err != nil {
		response.Fail(c, 403, err.Error())
		return
	}
	response.Success(c, gin.H{"list": list})
}

func (g *GroupApplyHandler) GetGroupApplyListByGroupId(c *gin.Context) {
	userId := c.GetInt64("id")
	groupIdStr := c.Param("group_id")
	groupId, err := strconv.ParseInt(groupIdStr, 10, 64)
	if err != nil {
		response.Fail(c, 400, "invalid group_id")
		return
	}
	list, err := g.groupApplyService.GetGroupApplyListByGroupId(c.Request.Context(), userId, groupId)
	if err != nil {
		response.Fail(c, 403, err.Error())
		return
	}
	response.Success(c, gin.H{"list": list})
}

func (g *GroupApplyHandler) GetGroupApplyDetail(c *gin.Context) {
	userId := c.GetInt64("id")
	applyIdStr := c.Param("apply_id")
	applyId, err := strconv.ParseInt(applyIdStr, 10, 64)
	if err != nil {
		response.Fail(c, 400, "invalid apply_id")
		return
	}
	apply, err := g.groupApplyService.GetGroupApplyDetail(c.Request.Context(), userId, applyId)
	if err != nil {
		response.Fail(c, 403, err.Error())
		return
	}
	response.Success(c, gin.H{"detail": apply})
}

func (g *GroupApplyHandler) DealGroupApply(c *gin.Context) {
	userId := c.GetInt64("id")
	groupIdStr := c.PostForm("group_id")
	groupId, _ := strconv.ParseInt(groupIdStr, 10, 64)
	var in req.DealWithGroupApplyReq
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, 400, "invalid request")
		return
	}
	err := g.groupApplyService.DealGroupApply(c.Request.Context(), userId, groupId, in.ApplyId, in.Status)
	if err != nil {
		response.Fail(c, 403, err.Error())
		return
	}
	response.Success(c, "处理成功")
}
func NewGroupApplyHandler(groupApplyService groupApplyService.GroupApplyServiceInterface) GroupApplyHandlerInterface {
	return &GroupApplyHandler{groupApplyService: groupApplyService}
}
