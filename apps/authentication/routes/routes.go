package routes

import (
	"server/apps/authentication/controllers"
	"server/foreox/settings"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

func GenRoutes(r *gin.Engine, subdir string) {
	r.POST(settings.API_PATH+subdir+"login", controllers.Login)
	r.POST(settings.API_PATH+subdir+"register", controllers.Register)
	r.GET(settings.API_PATH+subdir+"loggedin", middleware.Auth(), controllers.CheckLoggedIn)
}
