package routes

import (
	"server/apps/rooms/controllers"
	"server/foreox/settings"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

func GenRoutes(r *gin.Engine, subdir string) {
	r.POST(settings.API_PATH+subdir+"join/:id", middleware.Auth(), controllers.JoinRoom)
	r.GET(settings.API_PATH+subdir+"info/:id", middleware.Auth(), controllers.RoomInfo)
}
