package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	accounts "server/apps/accounts/models"
	"server/apps/messaging/services"
	"server/constants"
	"server/database"
	"server/errors"

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
	if client, ok := rm.CheckClientInRoom(user.(*accounts.User).Token); ok {
		client.Close()
	}

	services.ServeWs(rm, user.(*accounts.User), conn)
}

func NewRoom(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		fmt.Fprint(os.Stderr, errors.Wrap(err, err.Error()))
	}

	defer file.Close()

	// Reads the file and returns byte slice
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprint(os.Stderr, errors.Wrap(err, err.Error()))
	}

	// Uploading the file name
	uploadStream, err := database.Database.Bucket.OpenUploadStream(
		header.Filename,
	)
	if err != nil {
		fmt.Fprint(os.Stderr, errors.Wrap(err, err.Error()))
	}
	defer uploadStream.Close()

	// Writes the file to the database
	fileSize, err := uploadStream.Write(data)
	if err != nil {
		fmt.Fprint(os.Stderr, errors.Wrap(err, err.Error()))
	}

	log.Printf("Write file to DB was successful. File size: %d M\n", fileSize)

	var roomNameStruct roomName
	json.NewDecoder(c.Request.Body).Decode(&roomNameStruct)
	room := services.NewRoom(roomNameStruct.Name, header.Filename)
	fmt.Println(room)
	//TODO error
	user, _ := c.Get("user")
	go room.Serve()
	database.Database.Database.Collection(constants.USER_COLL).
		UpdateOne(context.Background(), bson.M{"_id": user.(*accounts.User).ID}, bson.M{"$push": bson.M{"rooms": room.Id}})
	c.JSON(http.StatusCreated, room)
}
