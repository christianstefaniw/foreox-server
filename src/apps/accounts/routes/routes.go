package routes

import (
	"server/src/apps/accounts/controllers"
	"server/src/devcord/settings"
	"server/src/middleware"

	"github.com/gin-gonic/gin"
)

func GenRoutes(r *gin.Engine, subdir string) {
	r.GET(settings.API_PATH+subdir+"info", middleware.Auth(), controllers.AccountInfo)
}
