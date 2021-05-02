package controllers

import (
	"context"
	"net/http"
	accounts "server/apps/accounts/models"
	rooms "server/apps/messaging/services"
	"server/constants"
	"server/database"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func JoinRoom(c *gin.Context) {
	rmId := c.Param("id")
	//TODO error
	user, _ := c.Get("user")
	database.Database.Database.Collection(constants.USER_COLL).
		UpdateOne(context.Background(), bson.M{"_id": user.(accounts.User).ID}, bson.M{"$push": bson.M{"Rooms": rmId}})
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

	for _, roomId := range user.(accounts.User).Rooms {
		var rm rooms.Room
		database.Database.FindOne(context.Background(), constants.ROOMS_COLL, bson.M{"_id": roomId}).Decode(&rm)
		allRooms = append(allRooms, rm)
	}
	c.JSON(http.StatusOK, allRooms)
}
