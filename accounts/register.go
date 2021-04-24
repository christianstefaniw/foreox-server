package accounts

import (
	"context"
	"errors"
	"server/constants"
	"server/database"
	"server/models"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// Registers user to the database
func Register(user *models.User) error {
	err := database.Collection.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(new(interface{}))
	if err.Error() != constants.MONGO_NO_DOC {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)

	if err != nil {
		return errors.New("error while hashing password")
	} else {
		user.Password = string(hash)
	}
	_, err = database.Collection.InsertOne(context.TODO(), user)
	if err != nil {
		return errors.New("error while creating user")
	} else {
		return nil
	}
}
