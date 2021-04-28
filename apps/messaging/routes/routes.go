package routes

import (
	"server/apps/messaging/controllers"

	"github.com/gin-gonic/gin"
)

func GenRoutes(r *gin.Engine) {
	r.GET("/ws/:id", controllers.ServeWs)
	r.GET("/api/newroom", controllers.NewRoom)
}
