package router

import (
	accounts "server/src/apps/accounts/routes"
	auth "server/src/apps/authentication/routes"
	home "server/src/apps/home/routes"
	rooms "server/src/apps/rooms/routes"
	"server/src/middleware"

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
