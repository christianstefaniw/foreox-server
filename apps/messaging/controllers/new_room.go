package controllers

import (
	"net/http"
	"server/apps/messaging/services"

	"github.com/gin-gonic/gin"
)

func NewRoom(c *gin.Context) {
	room := services.NewRoom()

	go room.Serve()

	c.Writer.WriteHeader(http.StatusCreated)
}
