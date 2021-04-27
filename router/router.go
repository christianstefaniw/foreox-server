package router

import (
	"server/controllers"
	"server/messaging"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

// Router is exported and used in main.go
func Router() *gin.Engine {

	router := gin.Default()
	rm := messaging.NewRoom()
	go rm.Serve()
	router.Use(middleware.CORSMiddleware())

	// Routes
	router.LoadHTMLFiles("index.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	router.POST("/api/user/login", controllers.LoginHandler)
	router.POST("/api/user/register", controllers.RegisterHandler)
	router.GET("/ws", func(c *gin.Context) {
		controllers.ServeWs(c, rm)
	})
	return router
}
