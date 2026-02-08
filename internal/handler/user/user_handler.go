package user

import (
	"net/http"

	"minichat/internal/req"
	"minichat/internal/service/user"

	"github.com/gin-gonic/gin"
)

var _ UserHandlerInterface = (*UserHandler)(nil)

type UserHandler struct {
	userService user.UserServiceInterface
}

func NewUserHandler(svc user.UserServiceInterface) *UserHandler {
	return &UserHandler{userService: svc}
}

func (h *UserHandler) Register(c *gin.Context) {
	var in req.RegisterReq
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "invalid request", "error": err.Error()})
		return
	}
	if err := h.userService.Register(c.Request.Context(), in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok"})
}

func (h *UserHandler) Login(c *gin.Context) {
	var in req.LoginReq
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "invalid request", "error": err.Error()})
		return
	}

	access, refresh, err := h.userService.Login(c.Request.Context(), in)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": gin.H{
			"access_token":  access,
			"refresh_token": refresh,
		},
	})
}

func (h *UserHandler) UpdateUserInfo(c *gin.Context) {
	var in req.UpdateUserReq
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "invalid request", "error": err.Error()})
	}
	userId, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "unauthorized"})
		return
	}
	if err := h.userService.UpdateUserInfo(c.Request.Context(), userId.(int64), in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok"})
}
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var in req.ChangePasswordReq
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "invalid request", "error": err.Error()})
	}
	userId, _ := c.Get("id")
	id := userId.(int64)
	if err := h.userService.ChangePassword(c, id, in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "密码修改失败", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "密码修改成功"})

}
func (h *UserHandler) ChangeUserId(c *gin.Context) {
	var in req.ChangeUserIdReq
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "invalid request", "error": err.Error()})
	}
	Id, _ := c.Get("id")
	id := Id.(int64)
	if err := h.userService.ChangeUserId(c, id, in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "userId修改失败", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "username修改成功"})
}

func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userId, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "unauthorized"})
		return
	}
	id := userId.(int64)
	userInfo, err := h.userService.GetUserInfo(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "获取用户信息失败", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": userInfo,
	})
}
func (h *UserHandler) CancelAccount(c *gin.Context) {
	var in req.CancelAccountReq
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "invalid request", "error": err.Error()})
		return
	}

	Id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "unauthorized"})
		return
	}
	id := Id.(int64)
	if err := h.userService.CancelAccount(c.Request.Context(), id, in.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "用户注销失败", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "用户注销成功"})
}
