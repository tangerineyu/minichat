package router

import "github.com/gin-gonic/gin"

func setupMiddleware(r *gin.Engine) {
	r.Use(gin.Recovery())

}
