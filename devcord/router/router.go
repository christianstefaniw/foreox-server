package router

import (
	accounts "server/apps/accounts/routes"
	auth "server/apps/authentication/routes"
	home "server/apps/home/routes"
	rooms "server/apps/rooms/routes"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

// Router is exported and used in main.go
func Router() *gin.Engine {

	router := gin.Default()
	router.Use(middleware.Cors())

	home.GenRoutes(router, "")
	auth.GenRoutes(router, "auth/")
	accounts.GenRoutes(router, "account/")
	rooms.GenRoutes(router, "rooms/")

	return router
}
