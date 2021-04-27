package controllers

import (
	"server/messaging"

	"github.com/gin-gonic/gin"
)

func NewRoom(c *gin.Context) {
	room := messaging.NewRoom()

	go room.Serve()

	c.Writer.WriteHeader(200)
}
