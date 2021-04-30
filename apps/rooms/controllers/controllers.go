package controllers

import (
	"context"
	accounts "server/apps/accounts/models"
	"server/database"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func JoinRoom(c *gin.Context) {
	rmId := c.Param("id")
	//TODO error
	user, _ := c.Get("user")

	database.Collection.UpdateOne(context.Background(), bson.M{"_id": user.(accounts.User).ID}, bson.M{"$push": bson.M{"Rooms": rmId}})
}

func RoomInfo(c *gin.Context) {
	rmId := c.Param("id")
	//TODO error
	user, _ := c.Get("user")

	database.Collection.UpdateOne(context.Background(), bson.M{"_id": user.(accounts.User).ID}, bson.M{"$push": bson.M{"Rooms": rmId}})
}
