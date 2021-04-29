package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"server/database"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type claims struct {
	Username string `json:"username"`
	Id       string `json:"id"`
	jwt.StandardClaims
}

var jwtKey = []byte("Foreox4224")

// Checks if users token is ok
func GetToken(c *gin.Context) {
	tknStr, err := c.Cookie("authToken")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			c.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Initialize a new instance of `Claims`
	claims := new(claims)

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		fmt.Println(err, "probably because you don't have your secrets the same")
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	userDocId, _ := primitive.ObjectIDFromHex(claims.Id)
	database.Collection.UpdateOne(context.Background(), bson.M{"_id": userDocId},
		bson.D{
			{Key: "$set", Value: bson.D{{Key: "token", Value: tknStr}}},
		})

	json.NewEncoder(c.Writer).Encode(claims)
}
