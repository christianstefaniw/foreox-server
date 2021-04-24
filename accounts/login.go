package accounts

import (
	"context"
	"fmt"
	"server/database"
	"server/models"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// Logs in user to the database
func Login(username, password string) (models.User, error) {
	var user models.User

	err := database.Collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		return user, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return user, err
	}

	fmt.Println(user.ID)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
	})

	tokenString, err := token.SignedString([]byte("Uploader4224"))

	if err != nil {
		return user, err
	}

	user.Token = tokenString
	user.Password = ""
	return user, nil
}
