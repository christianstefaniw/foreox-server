package routes

import (
	"server/apps/home/controllers"
	"server/foreox/settings"

	"github.com/gin-gonic/gin"
)

func GenRoutes(r *gin.Engine, subdir string) {
	r.GET(settings.API_PATH+subdir, controllers.Home)
}
