package controllers

import (
	"server/messaging"

	"github.com/gin-gonic/gin"
)

func ServeWs(r *messaging.Room, c *gin.Context) {
	messaging.ServeWs(r, c.Writer, c.Request)
}
