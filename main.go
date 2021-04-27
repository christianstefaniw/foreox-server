package main

import (
	"fmt"
	"log"
	"os"

	"server/database"
	"server/helpers"
	"server/middleware"
	"server/router"
)

func main() {
	initLogOutput()

	router := router.Router()

	helpers.LoadDotEnv()
	database.Connect()

	router.Use(middleware.CORSMiddleware())

	fmt.Println("Starting server on the port 8080...")
	router.Run(":" + os.Getenv("PORT"))
}

func initLogOutput() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)
}
