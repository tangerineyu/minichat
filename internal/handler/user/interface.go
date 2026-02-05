package user

import (
	"github.com/gin-gonic/gin"
)

type UserHandlerInterface interface {
	Register(c *gin.Context)

	Login(c *gin.Context)
	// 更改基础信息，昵称，头像
	UpdateUserInfo(c *gin.Context)
	// 更改密码
	ChangePassword(c *gin.Context)
	// 更改用户名
	ChangeUsername(c *gin.Context)
}
