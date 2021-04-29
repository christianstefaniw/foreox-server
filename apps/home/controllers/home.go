package controllers

import (
	"context"
	"encoding/json"
	accounts "server/apps/accounts/models"
	"server/database"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func Home(c *gin.Context) {
	// TODO handle error
	var user accounts.User
	tknStr, _ := c.Cookie("authToken")
	database.Collection.FindOne(context.Background(), bson.M{"token": tknStr}).Decode(&user)
	json.NewEncoder(c.Writer).Encode(user)
}
