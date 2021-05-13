package services

import (
	"context"
	accounts "server/apps/accounts/models"
	"server/constants"
	"server/database"
	errors "server/errors"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// Logs in user to the database
func Login(username, password string) (accounts.User, error) {
	var user accounts.User

	err := database.Database.FindOne(context.Background(), constants.USER_COLL, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err.Error() == constants.MONGO_NO_DOC {
			return user, errors.Error{Message: "Incorrect username or password"}
		}
		return user, errors.Wrap(err, err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return user, errors.Wrap(err, err.Error())
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
	})

	tokenString, err := token.SignedString([]byte("DevCord4224"))

	if err != nil {
		return user, errors.Wrap(err, err.Error())
	}

	user.Token = tokenString
	user.Password = ""
	return user, nil
}
