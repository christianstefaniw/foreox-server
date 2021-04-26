package router

import (
	"server/controllers"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

// Router is exported and used in main.go
func Router() *gin.Engine {

	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	// Routes
	router.POST("/api/user/login", controllers.LoginHandler)
	router.POST("/api/user/register", controllers.RegisterHandler)
	return router
}
