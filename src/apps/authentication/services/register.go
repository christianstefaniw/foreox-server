package services

import (
	"context"
	accounts "server/src/apps/accounts/models"
	"server/src/constants"
	"server/src/database"
	errors "server/src/errors"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// Registers user to the database
func Register(user *accounts.User) error {
	err := database.Database.FindOne(context.Background(), constants.USER_COLL, bson.M{"username": user.Username}).Decode(new(interface{}))
	if err.Error() != constants.MONGO_NO_DOC {
		return errors.Wrap(err, err.Error())
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)

	if err != nil {
		return errors.Wrap(err, "error hashing password")
	}

	user.Password = string(hash)

	_, err = database.Database.InsertOne(context.Background(), constants.USER_COLL, user)
	if err != nil {
		return errors.Wrap(err, "error while creating user")
	}
	return nil
}
