package friend_apply

import (
	"minichat/internal/handler/response"

	"minichat/internal/req"
	"minichat/internal/service/friend_apply"

	"github.com/gin-gonic/gin"
)

var _ FriendApplyHandlerInterface = (*FriendApplyHandler)(nil)

type FriendApplyHandler struct {
	friendApplyService friend_apply.FriendApplyServiceInterface
}

func (h *FriendApplyHandler) GetFriendApplyList(c *gin.Context) {
	id := c.GetInt64("id")
	list, err := h.friendApplyService.GetFriendApply(c.Request.Context(), id)
	if err != nil {
		response.ServerError(c, err)
		return
	}
	response.Success(c, gin.H{"list": list})
}

func (h *FriendApplyHandler) DealWithFriendApply(c *gin.Context) {
	var in req.DealWithFriendApplyReq
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, 400, "invalid request")
		return
	}
	Id := c.MustGet("id")

	if err := h.friendApplyService.DealWithFriendApply(c.Request.Context(), Id.(int64), in); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.Success(c, nil)
}

func NewUserHandler(svc friend_apply.FriendApplyServiceInterface) *FriendApplyHandler {
	return &FriendApplyHandler{
		friendApplyService: svc,
	}
}

func (h *FriendApplyHandler) SendFriendApply(c *gin.Context) {
	var in req.SendFriendApplyReq
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, 400, "invalid request")
		return
	}

	fromUserId, exists := c.Get("id")
	if !exists {
		response.Fail(c, 401, "unauthorized")
		return
	}

	if err := h.friendApplyService.SendFriendApply(c.Request.Context(), fromUserId.(int64), in); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}

	response.Success(c, nil)
}
