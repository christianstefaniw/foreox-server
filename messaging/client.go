package messaging

import (
	"context"
	"fmt"
	errors "server/errors"
	"server/models"
	"time"

	"github.com/gorilla/websocket"
)

const (
	readDeadline  = 5 * time.Second
	writeDeadline = 5 * time.Second
)

// TODO handle errors
// TODO make msg chan buffered and handle when buffer gets too big
type client struct {
	room   *Room
	conn   *websocket.Conn
	msg    chan []byte
	ctx    context.Context
	cancel context.CancelFunc
	models.User
}

func (c *client) unregister() {
	c.room.unregister <- c
	c.conn.Close()
}

func (c *client) read(deadline time.Duration) {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				errors.PrintError(errors.GetErrorKey(), errors.Wrap(err, err.Error()))
			}
			c.cancel()
		} else {
			c.room.broadcast <- message
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
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *client) writeSteward() {
	const timeout = 5 * time.Second
	heartbeat := make(chan interface{})

	go c.write(heartbeat, timeout/2)

	for {
		select {
		case <-heartbeat:
		case <-c.ctx.Done():
			close(heartbeat)
			return
		case <-time.After(timeout):
			fmt.Println("writer unhealthy, unregistering...")
			c.cancel()
		}
	}
}

func (c *client) doWork() {

	go c.read(readDeadline)
	go c.writeSteward()

	for {
		select {
		case <-c.ctx.Done():
			fmt.Println("client closed")
			c.unregister()
			return

		case <-c.room.ctx.Done():
			c.unregister()
			return
		}
	}
}

func ServeWs(r *Room, conn *websocket.Conn) {
	ctx, cancel := context.WithCancel(context.Background())
	c := &client{room: r, conn: conn, msg: make(chan []byte), ctx: ctx, cancel: cancel}

	c.room.register <- c

	go c.doWork()
}
