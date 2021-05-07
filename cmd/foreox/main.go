package main

import (
	"context"
	"log"
	"os"

	"server/database"
	"server/foreox/router"
	"server/foreox/settings"
	"server/helpers"
)

func main() {
	initLogOutput()
	helpers.LoadDotEnv()

	database.Connect()
	defer database.Database.Database.Client().Disconnect(context.Background())

	settings.Settings()

	router := router.Router()

	router.Run(":" + os.Getenv("PORT"))
}

func initLogOutput() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)
}
