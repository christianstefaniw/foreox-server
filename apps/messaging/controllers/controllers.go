package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	accounts "server/apps/accounts/models"
	"server/apps/messaging/services"
	"server/constants"
	"server/database"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
)

type roomName struct {
	Name string `json:"roomName"`
}

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
	user, _ := c.Get("user")
	if rm.CheckClientInRoom(user.(*accounts.User).Token) {
		return
	}
	services.ServeWs(rm, user.(*accounts.User), conn)
}

func NewRoom(c *gin.Context) {
	var roomNameStruct roomName
	json.NewDecoder(c.Request.Body).Decode(&roomNameStruct)
	room := services.NewRoom(roomNameStruct.Name)
	//TODO error
	user, _ := c.Get("user")
	go room.Serve()
	database.Database.Database.Collection(constants.USER_COLL).
		UpdateOne(context.Background(), bson.M{"_id": user.(*accounts.User).ID}, bson.M{"$push": bson.M{"rooms": room.Id}})
	c.JSON(http.StatusCreated, room)
}
