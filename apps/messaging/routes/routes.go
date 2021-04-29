package routes

import (
	"server/apps/messaging/controllers"
	"server/foreox/settings"

	"github.com/gin-gonic/gin"
)

func GenRoutes(r *gin.Engine, subdir string) {
	r.GET(settings.API_PATH+subdir+"/ws/:id", controllers.ServeWs)
	r.GET(settings.API_PATH+subdir+"/newroom", controllers.NewRoom)
}
