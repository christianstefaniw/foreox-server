package routes

import (
	"server/src/apps/home/controllers"
	"server/src/devcord/settings"

	"github.com/gin-gonic/gin"
)

func GenRoutes(r *gin.Engine, subdir string) {
	r.GET(settings.API_PATH+subdir, controllers.Home)
}
