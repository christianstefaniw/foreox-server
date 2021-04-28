package routes

import (
	"server/apps/authentication/controllers"

	"github.com/gin-gonic/gin"
)

func GenRoutes(r *gin.Engine) {
	r.POST("/api/auth/login", controllers.Login)
	r.POST("/api/auth/register", controllers.Register)
	r.GET("/api/auth/gettoken", controllers.GetToken)
}
