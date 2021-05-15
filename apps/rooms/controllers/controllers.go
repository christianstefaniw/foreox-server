package controllers

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	accounts "server/apps/accounts/models"
	messaging "server/apps/messaging/services"
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

	room := messaging.NewRoom(c.Request.FormValue("roomName"), header.Filename)
	fmt.Println(c.Request.FormValue("roomName"))
	//TODO error
	user, _ := c.Get("user")
	go room.Serve()
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
	path := helpers.RootDir() + constants.MEDIA_DIR + fileName
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
	// getting all images in gridfs
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var results bson.M
	err := database.Database.FindOne(ctx, constants.FILES_COLL, bson.M{}).Decode(&results)
	if err != nil {
		fmt.Println(err)
	}
	// gridfs downloading the file to cmd/foreox
	var buf bytes.Buffer
	for i := range allRooms {
		dStream, err := database.Database.Bucket.DownloadToStreamByName(allRooms[i].Image, &buf)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(dStream)
			ioutil.WriteFile(allRooms[i].Image, buf.Bytes(), 0600)
		}
	}

	c.JSON(http.StatusOK, allRooms)
}
