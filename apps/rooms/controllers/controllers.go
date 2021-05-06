package controllers

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	accounts "server/apps/accounts/models"
	rooms "server/apps/messaging/services"
	"server/constants"
	"server/database"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func JoinRoom(c *gin.Context) {
	rmId := c.Param("id")
	//TODO error
	user, _ := c.Get("user")
	database.Database.Database.Collection(constants.USER_COLL).
		UpdateOne(context.Background(), bson.M{"_id": user.(*accounts.User).ID}, bson.M{"$push": bson.M{"Rooms": rmId}})
}

func RoomInfo(c *gin.Context) {
	var rm rooms.Room
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
	var allRooms []rooms.Room
	user, _ := c.Get("user")
	for _, roomId := range user.(*accounts.User).Rooms {
		var rm rooms.Room
		database.Database.FindOne(context.Background(), constants.ROOMS_COLL, bson.M{"_id": roomId}).Decode(&rm)
		allRooms = append(allRooms, rm)
	}
	// getting all images in gridfs
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    var results bson.M
	err := database.Database.FindOne(ctx, constants.FILES_COLL, bson.M{}).Decode(&results)
	if err != nil {
		fmt.Println(err)
	}
	// gridfs downloading the file to cmd/foreox
	var buf bytes.Buffer
	for i, _ := range allRooms {
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