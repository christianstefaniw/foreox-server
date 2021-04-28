package routes

import (
	"server/apps/home/controllers"

	"github.com/gin-gonic/gin"
)

func GenRoutes(r *gin.Engine) {
	r.GET("/api/", controllers.Home)
}
