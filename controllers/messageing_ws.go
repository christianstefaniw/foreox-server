package controllers

import (
	"server/messaging"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// specifies parameters for upgrading an http connection to a ws connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ServeWs(c *gin.Context, rm *messaging.Room) {
	conn, _ := upgrader.Upgrade(c.Writer, c.Request, nil)
	messaging.ServeWs(rm, conn)
}
