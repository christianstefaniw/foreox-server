package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	accounts "server/apps/accounts/models"
	messaging "server/apps/messaging/services"
	"server/apps/rooms/services"
	"server/constants"
	"server/database"
	"server/errors"
	"server/helpers"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewRoom(c *gin.Context) {
	c.Request.ParseForm()
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		fmt.Fprint(os.Stderr, errors.Wrap(err, err.Error()))
	}
	defer file.Close()

	fileName, err := services.SaveImage(file, header)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	room := messaging.NewRoom(c.Request.FormValue("roomName"), fileName)

	go room.Serve()

	//TODO error
	user, _ := c.Get("user")
	database.Database.Database.Collection(constants.USER_COLL).
		UpdateOne(context.Background(), bson.M{"_id": user.(*accounts.User).ID}, bson.M{"$push": bson.M{"rooms": room.Id}})

	c.JSON(http.StatusCreated, room)
}

func ServeRoom(c *gin.Context) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	var rm *messaging.Room
	rmId := c.Param("id")

	rm, ok := messaging.GetRoom(rmId)
	if !ok {
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}

	conn, _ := upgrader.Upgrade(c.Writer, c.Request, nil)
	user, _ := c.Get("user")
	if client, ok := rm.CheckClientInRoom(user.(*accounts.User).Token); ok {
		client.Close()
	}

	messaging.ServeWs(rm, user.(*accounts.User), conn)
}

func ServeRoomImage(c *gin.Context) {
	fileName := c.Param("name")
	path := fmt.Sprintf("%s/%s/%s", helpers.RootDir(), constants.MEDIA_DIR, fileName)
	//TODO error
	file, _ := os.Open(path)
	http.ServeContent(c.Writer, c.Request, "room_image", time.Now(), file)
}

func JoinRoom(c *gin.Context) {
	rmId := c.Param("id")
	//TODO error
	user, _ := c.Get("user")
	database.Database.Database.Collection(constants.USER_COLL).
		UpdateOne(context.Background(), bson.M{"_id": user.(*accounts.User).ID}, bson.M{"$push": bson.M{"Rooms": rmId}})
}

func RoomInfo(c *gin.Context) {
	var rm messaging.Room
	//TODO error
	rmId, _ := primitive.ObjectIDFromHex(c.Param("id"))
	//TODO error
	user, _ := c.Get("user")
	for _, userRoomId := range user.(*accounts.User).Rooms {
		if userRoomId == rmId {
			database.Database.FindOne(context.Background(), constants.ROOMS_COLL, bson.M{"_id": rmId}).Decode(&rm)
		}
	}

	if rm.Id == (primitive.ObjectID{}) {
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, rm)
}

func AllUsersRooms(c *gin.Context) {
	var allRooms []messaging.Room
	user, _ := c.Get("user")
	for _, roomId := range user.(*accounts.User).Rooms {
		var rm messaging.Room
		database.Database.FindOne(context.Background(), constants.ROOMS_COLL, bson.M{"_id": roomId}).Decode(&rm)
		allRooms = append(allRooms, rm)
	}

	c.JSON(http.StatusOK, allRooms)
}
