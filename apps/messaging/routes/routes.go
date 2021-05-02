package routes

import (
	"server/apps/messaging/controllers"
	"server/foreox/settings"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

func GenRoutes(r *gin.Engine, subdir string) {
	r.GET(settings.API_PATH+subdir+"ws/:id", middleware.Auth(), controllers.ServeWs)
	r.POST(settings.API_PATH+subdir+"newroom", middleware.Auth(), controllers.NewRoom)
}
