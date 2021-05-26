package routes

import (
	"server/src/apps/authentication/controllers"
	"server/src/devcord/settings"
	"server/src/middleware"

	"github.com/gin-gonic/gin"
)

func GenRoutes(r *gin.Engine, subdir string) {
	r.POST(settings.API_PATH+subdir+"login", controllers.Login)
	r.POST(settings.API_PATH+subdir+"register", controllers.Register)
	r.GET(settings.API_PATH+subdir+"loggedin", middleware.Auth(), controllers.CheckLoggedIn)
}
