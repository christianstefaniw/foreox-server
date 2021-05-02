package middleware

import (
	"context"
	"fmt"
	"net/http"
	accounts "server/apps/accounts/models"
	"server/constants"
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

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tknStr, err := c.Cookie("authToken")
		if err != nil {
			if err == http.ErrNoCookie {
				// If the cookie is not set, return an unauthorized status
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			// For any other type of error, return a bad request status
			c.AbortWithStatus(http.StatusBadRequest)
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
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			fmt.Println(err, "probably because you don't have your secrets the same")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		var user accounts.User
		userDocId, _ := primitive.ObjectIDFromHex(claims.Id)
		userDoc := database.Database.FindOneAndUpdate(context.Background(), constants.USER_COLL, bson.M{"_id": userDocId},
			bson.D{
				{Key: "$set", Value: bson.D{{Key: "token", Value: tknStr}}},
			})
		if userDoc.Err() != nil {
			if userDoc.Err().Error() == constants.MONGO_NO_DOC {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		userDoc.Decode(&user)
		c.Set("user", &user)
	}
}
