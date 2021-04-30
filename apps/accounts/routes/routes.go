package routes

import (
	"server/apps/accounts/controllers"
	"server/foreox/settings"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

func GenRoutes(r *gin.Engine, subdir string) {
	r.GET(settings.API_PATH+subdir+"/info", middleware.Auth(), controllers.AccountInfo)
}
