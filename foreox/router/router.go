package router

import (
	auth "server/apps/authentication/routes"
	home "server/apps/home/routes"
	messaging "server/apps/messaging/routes"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

// Router is exported and used in main.go
func Router() *gin.Engine {

	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	home.GenRoutes(router)
	auth.GenRoutes(router)
	messaging.GenRoutes(router)

	// Routes
	// router.LoadHTMLFiles("index.html")
	// router.GET("/", func(c *gin.Context) {
	// 	c.HTML(200, "index.html", nil)
	// })
	return router
}
