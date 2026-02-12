package user

import (
	"minichat/internal/handler/response"
	"minichat/internal/req"
	"minichat/internal/service/user"
	"net/http"

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
		response.Fail(c, 400, "invalid request")
		return
	}
	if err := h.userService.Register(c.Request.Context(), in); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *UserHandler) Login(c *gin.Context) {
	var in req.LoginReq
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, 400, "invalid request")
		return
	}

	access, refresh, err := h.userService.Login(c.Request.Context(), in)
	if err != nil {
		response.Fail(c, 401, err.Error())
		return
	}

	response.Success(c, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func (h *UserHandler) UpdateUserInfo(c *gin.Context) {
	var in req.UpdateUserReq
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, 400, "invalid request")
		return
	}
	userId, exists := c.Get("id")
	if !exists {
		response.Fail(c, 401, "unauthorized")
		return
	}
	if err := h.userService.UpdateUserInfo(c.Request.Context(), userId.(int64), in); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	var in req.ChangePasswordReq
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, 400, "invalid request")
		return
	}
	userId, exists := c.Get("id")
	if !exists {
		response.Fail(c, 401, "unauthorized")
		return
	}
	id := userId.(int64)
	if err := h.userService.ChangePassword(c, id, in); err != nil {
		response.Fail(c, 400, "密码修改失败")
		return
	}
	response.Success(c, gin.H{"msg": "密码修改成功"})
}

func (h *UserHandler) ChangeUserId(c *gin.Context) {
	var in req.ChangeUserIdReq
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, 400, "invalid request")
		return
	}
	Id, exists := c.Get("id")
	if !exists {
		response.Fail(c, 401, "unauthorized")
		return
	}
	id := Id.(int64)
	if err := h.userService.ChangeUserId(c, id, in); err != nil {
		response.Fail(c, 400, "userId修改失败")
		return
	}
	response.Success(c, gin.H{"msg": "userId修改成功"})
}

func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userId, exists := c.Get("id")
	if !exists {
		response.Fail(c, 401, "unauthorized")
		return
	}
	id := userId.(int64)
	userInfo, err := h.userService.GetUserInfo(c, id)
	if err != nil {
		response.Fail(c, 400, "获取用户信息失败")
		return
	}
	response.Success(c, userInfo)
}

func (h *UserHandler) CancelAccount(c *gin.Context) {
	var in req.CancelAccountReq
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, 400, "invalid request")
		return
	}

	Id, exists := c.Get("id")
	if !exists {
		response.Fail(c, 401, "unauthorized")
		return
	}
	id := Id.(int64)
	if err := h.userService.CancelAccount(c.Request.Context(), id, in.Password); err != nil {
		response.Fail(c, 400, "用户注销失败")
		return
	}
	response.Success(c, gin.H{"msg": "用户注销成功"})
}

func (h *UserHandler) UploadAvatar(c *gin.Context) {
	id, exists := c.Get("id")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	userID := id.(int64)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "missing file")
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		// 兜底：有些客户端可能不带
		contentType = c.GetHeader("Content-Type")
	}

	avatarURL, err := h.userService.UploadAvatar(c.Request.Context(), userID, contentType, file)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{"avatar_url": avatarURL})
}
