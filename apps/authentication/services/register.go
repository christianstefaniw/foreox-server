package services

import (
	"context"
	accounts "server/apps/accounts/models"
	"server/constants"
	"server/database"
	errors "server/errors"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// Registers user to the database
func Register(user *accounts.User) error {
	err := database.Collection.FindOne(context.Background(), bson.M{"username": user.Username}).Decode(new(interface{}))
	if err.Error() != constants.MONGO_NO_DOC {
		return errors.Wrap(err, err.Error())
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)

	if err != nil {
		return errors.Wrap(err, "error hashing password")
	}

	user.Password = string(hash)

	_, err = database.Collection.InsertOne(context.Background(), user)
	if err != nil {
		return errors.Wrap(err, "error while creating user")
	}
	return nil
}
