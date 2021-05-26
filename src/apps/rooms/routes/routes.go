package routes

import (
	"server/src/apps/rooms/controllers"
	"server/src/devcord/settings"
	"server/src/middleware"

	"github.com/gin-gonic/gin"
)

func GenRoutes(r *gin.Engine, subdir string) {
	r.POST(settings.API_PATH+subdir+"join/:id", middleware.Auth(), controllers.JoinRoom)
	r.GET(settings.API_PATH+subdir+"info/:id", middleware.Auth(), controllers.RoomInfo)
	r.GET(settings.API_PATH+subdir+"allrooms", middleware.Auth(), controllers.AllUsersRooms)
	r.GET(settings.API_PATH+subdir+"img/:name", controllers.ServeRoomImage)
	r.GET(settings.API_PATH+subdir+"room/:id", middleware.Auth(), controllers.ServeRoom)
	r.POST(settings.API_PATH+subdir+"newroom", middleware.Auth(), controllers.NewRoom)
}
