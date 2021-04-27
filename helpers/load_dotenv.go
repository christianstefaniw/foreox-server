package helpers

import (
	"log"
	errors "server/errors"

	"github.com/joho/godotenv"
)

// load environment
func LoadDotEnv() {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(errors.Wrap(err, err.Error()))
	}
}
