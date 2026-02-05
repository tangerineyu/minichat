package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int    `json:"code"`
	Data any    `json:"data,omitempty"`
	Msg  string `json:"msg"`
	Err  error  `json:"-"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: 0, Msg: "ok", Data: data})
}

func Fail(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusBadRequest, Response{Code: code, Msg: msg})
}
