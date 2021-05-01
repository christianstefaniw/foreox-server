package main

import (
	"context"
	"fmt"
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

	fmt.Println("Starting server on the port 8080...")
	router.Run(":" + os.Getenv("PORT"))
}

func initLogOutput() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)
}
