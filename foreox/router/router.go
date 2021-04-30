package router

import (
	accounts "server/apps/accounts/routes"
	auth "server/apps/authentication/routes"
	home "server/apps/home/routes"
	messaging "server/apps/messaging/routes"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

// Router is exported and used in main.go
func Router() *gin.Engine {

	router := gin.Default()
	router.Use(middleware.Cors())

	home.GenRoutes(router, "")
	auth.GenRoutes(router, "auth")
	messaging.GenRoutes(router, "messaging")
	accounts.GenRoutes(router, "account")

	// Routes
	// router.LoadHTMLFiles("index.html")
	// router.GET("/", func(c *gin.Context) {
	// 	c.HTML(200, "index.html", nil)
	// })
	return router
}
