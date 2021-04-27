package messaging

import (
	"fmt"
	"server/errors"
	"server/models"
	"time"

	"github.com/gorilla/websocket"
)

// TODO handle errors
// TODO make msg chan buffered and handle when buffer gets too big
type client struct {
	room *Room
	conn *websocket.Conn
	msg  chan []byte
	done chan interface{}
	models.User
}

func (c *client) unregister() {
	close(c.done)
	c.room.unregister <- c
	c.conn.Close()
}

func (c *client) read(heartbeat chan interface{}, pulseInterval time.Duration) {
	pulse := time.NewTicker(pulseInterval)
	message := make(chan []byte)

	go func() {
		// read message sent to THIS connection
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if err.Error() != "websocket: close 1001 (goig away)" {
				errors.PrintError(errors.GetErrorKey(), errors.Wrap(err, err.Error()))
			}
			fmt.Println("disconnecting...")
			c.unregister()
			return
		}
		message <- msg
	}()

	for {
		select {
		case <-pulse.C:
			sendPulse(heartbeat)
		case msg := <-message:
			// send message to all clients in room
			c.room.broadcast <- msg
		case _, ok := <-c.done:
			if !ok {
				return
			}
		}
	}
}

func (c *client) write(heartbeat chan interface{}, pulseInterval time.Duration) {
	pulse := time.NewTicker(pulseInterval)
	for {
		select {
		case <-pulse.C:
			sendPulse(heartbeat)
		case msg := <-c.msg:
			c.conn.WriteMessage(websocket.TextMessage, msg)
		case _, ok := <-c.done:
			if !ok {
				return
			}
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
		case _, ok := <-c.done:
			fmt.Println("done channel closed")
			if !ok {
				return
			}
		case <-c.room.ctx.Done():
			c.unregister()
			return
		case <-time.After(timeout):
			c.unregister()
			return
		}
	}
}

func ServeWs(r *Room, conn *websocket.Conn) {
	c := &client{room: r, conn: conn, msg: make(chan []byte), done: make(chan interface{})}

	c.room.register <- c

	go c.doWork()
}
