package messaging

import (
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO make msg chan buffered and handle when buffer gets too big
type client struct {
	id   primitive.ObjectID
	room *room
	conn websocket.Conn
	msg  chan []byte
}

// specifies parameters for upgrading an http connection to a ws connection
var upgrader = new(websocket.Upgrader)

func (c *client) unregister() {
	c.room.unregister <- c
	c.conn.Close()
}

func (c *client) read() {
	defer c.unregister()

}
