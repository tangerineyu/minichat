package ws

import "github.com/gin-gonic/gin"

type WSHandlerInterface interface {
	// Connect upgrades an HTTP request to a WebSocket connection and registers the client to the hub.
	Connect(c *gin.Context)
}
