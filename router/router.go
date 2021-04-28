package router

import (
	"fmt"
	"server/controllers"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

// Router is exported and used in main.go
func Router() *gin.Engine {

	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	// Routes
	// router.LoadHTMLFiles("index.html")
	// router.GET("/", func(c *gin.Context) {
	// 	c.HTML(200, "index.html", nil)
	// })
	router.GET("/", func(c *gin.Context) {
		fmt.Fprint(c.Writer, "Welcome!")
	})
	router.POST("/api/user/login", controllers.LoginHandler)
	router.POST("/api/user/register", controllers.RegisterHandler)
	router.GET("/api/newroom", controllers.NewRoom)
	router.GET("/api/getToken", controllers.GetToken)
	router.GET("/ws/:id", controllers.ServeWs)
	return router
}
