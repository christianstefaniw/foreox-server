package messaging

import (
	"net/http"
	"server/models"

	"github.com/gorilla/websocket"
)

// TODO handle errors
// TODO make msg chan buffered and handle when buffer gets too big
type client struct {
	room *Room
	conn *websocket.Conn
	msg  chan []byte
	models.User
}

// specifies parameters for upgrading an http connection to a ws connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (c *client) unregister() {
	c.room.unregister <- c
	c.conn.Close()
}

func (c *client) closeConn() {
	c.conn.Close()
}

func (c *client) read() {
	defer c.closeConn()
	for {
		// read message sent to THIS connection
		_, msg, _ := c.conn.ReadMessage()
		// send message to all clients in room
		c.room.broadcast <- msg
	}
}

func (c *client) write() {
	defer c.unregister()
	for {
		msg := <-c.msg
		c.conn.WriteMessage(websocket.TextMessage, msg)
	}
}

func ServeWs(r *Room, w http.ResponseWriter, req *http.Request) {
	serveWs(r, w, req)
}

func serveWs(r *Room, w http.ResponseWriter, req *http.Request) {
	conn, _ := upgrader.Upgrade(w, req, nil)
	c := &client{room: r, conn: conn, msg: make(chan []byte)}

	c.room.register <- c

	go c.write()
	go c.read()
}
