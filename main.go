package main

import (
	"fmt"
	"os"

	"server/database"
	"server/helpers"
	"server/middleware"
	"server/router"
)

func main() {
	router := router.Router()

	helpers.LoadDotEnv()
	database.Connect()

	router.Use(middleware.CORSMiddleware())

	fmt.Println("Starting server on the port 8080...")
	router.Run(os.Getenv("PORT"))
}
