package friend_apply

import (
	"net/http"

	"minichat/internal/req"
	"minichat/internal/service/friend_apply"

	"github.com/gin-gonic/gin"
)

var _ FriendApplyHandlerInterface = (*FriendApplyHandler)(nil)

type FriendApplyHandler struct {
	friendApplyService friend_apply.FriendApplyServiceInterface
}

func (h *FriendApplyHandler) DealWithFriendApply(c *gin.Context) {
	var in req.DealWithFriendApplyReq
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "invalid request", "error": err.Error()})
		return
	}
	Id := c.MustGet("id")

	if err := h.friendApplyService.DealWithFriendApply(c.Request.Context(), Id.(int64), in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok"})
}

func NewUserHandler(svc friend_apply.FriendApplyServiceInterface) *FriendApplyHandler {
	return &FriendApplyHandler{
		friendApplyService: svc,
	}
}

func (h *FriendApplyHandler) SendFriendApply(c *gin.Context) {
	var in req.SendFriendApplyReq
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "invalid request", "error": err.Error()})
		return
	}

	// JWT middleware stores user id under key "id"
	fromUserId, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "unauthorized"})
		return
	}

	if err := h.friendApplyService.SendFriendApply(c.Request.Context(), fromUserId.(int64), in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok"})
}
