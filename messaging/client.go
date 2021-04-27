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
	readWait   = 5 * time.Second
	writeWait  = 5 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
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
	close(c.msg)
	c.room.unregister <- c
	c.conn.Close()
}

func (c *client) read() {
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				errors.PrintError(errors.GetErrorKey(), errors.Wrap(err, err.Error()))
			}
			return
		} else {
			c.room.broadcast <- message
		}
	}
}

func (c *client) write() {
	ticker := time.NewTicker(pingPeriod)

	for {
		select {
		case msg := <-c.msg:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			c.conn.WriteMessage(websocket.TextMessage, msg)
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.cancel()
				return
			}
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *client) doWork() {

	go c.read()
	go c.write()

	for {
		select {
		case <-c.ctx.Done():
			c.unregister()
			fmt.Println("client closed")
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
