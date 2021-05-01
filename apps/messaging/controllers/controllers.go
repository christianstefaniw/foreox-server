package controllers

import (
	"context"
	"net/http"
	accounts "server/apps/accounts/models"
	"server/apps/messaging/services"
	"server/constants"
	"server/database"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ServeWs(c *gin.Context) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	var rm *services.Room
	rmId := c.Param("id")

	rm, ok := services.GetRoom(rmId)
	if !ok {
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}

	conn, _ := upgrader.Upgrade(c.Writer, c.Request, nil)
	services.ServeWs(rm, conn)
}

func NewRoom(c *gin.Context) {
	room := services.NewRoom()
	//TODO error
	user, _ := c.Get("user")
	go room.Serve()
	database.Database.Database.Collection(constants.USER_COLL).
		UpdateOne(context.Background(), bson.M{"_id": user.(accounts.User).ID}, bson.M{"$push": bson.M{"rooms": room.Id}})
	c.JSON(http.StatusCreated, room)
}
