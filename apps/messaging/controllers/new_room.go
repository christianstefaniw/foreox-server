package controllers

import (
	"net/http"
	"server/apps/messaging/services"

	"github.com/gin-gonic/gin"
)

type roomId struct {
	Id string `json:"id"`
}

func NewRoom(c *gin.Context) {
	room := services.NewRoom()
	rmIdStream := make(chan string, 1)

	go room.Serve(rmIdStream)

	roomId := roomId{<-rmIdStream}

	c.JSON(http.StatusCreated, roomId)
}
