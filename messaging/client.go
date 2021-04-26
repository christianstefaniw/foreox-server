package messaging

import (
	"context"
	"fmt"
	"server/models"
	"time"

	"github.com/gorilla/websocket"
)

// TODO handle errors
// TODO make msg chan buffered and handle when buffer gets too big
// TODO make client pool, I think that might work well. TBD tho
type client struct {
	room *Room
	conn *websocket.Conn
	msg  chan []byte
	ctx  context.Context
	models.User
}

func (c *client) unregister() {
	c.room.unregister <- c
	c.conn.Close()
}

func (c *client) read(heartbeat chan interface{}, pulseInterval time.Duration) {
	defer c.unregister()
	pulse := time.NewTicker(pulseInterval)
	message := make(chan []byte)

	go func() {
		// read message sent to THIS connection
		_, msg, _ := c.conn.ReadMessage()
		message <- msg
	}()

	for {
		select {
		case <-pulse.C:
			sendPulse(heartbeat)
		case msg := <-message:
			// send message to all clients in room
			c.room.broadcast <- msg
		}

	}
}

func (c *client) write(heartbeat chan interface{}, pulseInterval time.Duration) {
	defer c.unregister()
	pulse := time.NewTicker(pulseInterval)
	for {
		select {
		case <-pulse.C:
			sendPulse(heartbeat)
		case msg := <-c.msg:
			c.conn.WriteMessage(websocket.TextMessage, msg)
		}
	}
}

func (c *client) doWork() {
	writeHeartbeat := make(chan interface{}, 1)
	readHeartbeat := make(chan interface{}, 1)
	const timeout = 10 * time.Second

	go c.write(writeHeartbeat, timeout/2)
	go c.read(readHeartbeat, timeout/2)

	for {
		select {
		case _, ok := <-readHeartbeat:
			fmt.Println("heartbeat read")
			if !ok {
				return
			}
		case _, ok := <-writeHeartbeat:
			fmt.Println("heartbeat write")
			if !ok {
				return
			}
		case <-c.ctx.Done():
			c.unregister()
			return
		case <-time.After(timeout):
			c.unregister()
			return
		}
	}
}

func ServeWs(r *Room, conn *websocket.Conn) {
	c := &client{room: r, conn: conn, msg: make(chan []byte), ctx: r.ctx}

	c.room.register <- c

	go c.doWork()
}
