package middlerware

import (
	"github.com/gin-gonic/gin"
	jwt_util "minichat/util/jwt"
	"strings"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			//拦截请求，不再向下
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header is missing"})
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "token格式错误"})
			c.Abort()
			return
		}
		claims, err := jwt_util.ValidateAccessToken(parts[1])
		if err != nil {
			c.JSON(401, gin.H{"error": "无效或已过期的token"})
			c.Abort()
			return
		}
		c.Set("id", claims.Id)
		c.Set("username", claims.Username)
		c.Next()
	}
}
