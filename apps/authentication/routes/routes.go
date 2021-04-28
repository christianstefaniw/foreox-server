package routes

import (
	"server/apps/authentication/controllers"
	"server/foreox/settings"

	"github.com/gin-gonic/gin"
)

func GenRoutes(r *gin.Engine, subdir string) {
	r.POST(settings.API_PATH+subdir+"/login", controllers.Login)
	r.POST(settings.API_PATH+subdir+"/register", controllers.Register)
	r.GET(settings.API_PATH+subdir+"/gettoken", controllers.GetToken)
}
