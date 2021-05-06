package helpers

import (
	"log"
	errors "server/errors"

	"github.com/joho/godotenv"
)

func LoadDotEnv() {
	path := RootDir() + ".env"

	err := godotenv.Load(path)

	if err != nil {
		log.Fatal(errors.Wrap(err, err.Error()))
	}
}
