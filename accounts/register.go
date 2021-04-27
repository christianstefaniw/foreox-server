package accounts

import (
	"context"
	"server/constants"
	"server/database"
	errors "server/errors"
	"server/models"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// Registers user to the database
func Register(user *models.User) error {
	err := database.Collection.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(new(interface{}))
	if err.Error() != constants.MONGO_NO_DOC {
		return errors.Wrap(err, err.Error())
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)

	if err != nil {
		return errors.Wrap(err, "error hashing password")
	} else {
		user.Password = string(hash)
	}
	_, err = database.Collection.InsertOne(context.TODO(), user)
	if err != nil {
		return errors.Wrap(err, "error while creating user")
	} else {
		return nil
	}
}
