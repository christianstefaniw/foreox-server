package services

import (
	"context"
	"fmt"
	accounts "server/apps/accounts/models"
	"server/errors"
	"time"

	"github.com/gorilla/websocket"
)

type client struct {
	room   *Room
	conn   *websocket.Conn
	msg    chan []byte
	ctx    context.Context
	cancel context.CancelFunc
	*accounts.User
}

const (
	readWait       = 5 * time.Second
	writeWait      = 5 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newLine = []byte{'\n'}
)

func ServeWs(r *Room, user *accounts.User, conn *websocket.Conn) {
	ctx, cancel := context.WithCancel(context.Background())
	c := &client{room: r, conn: conn, msg: make(chan []byte, 256), ctx: ctx, cancel: cancel, User: user}
	c.room.register <- c
	go c.doWork()
}

func (c *client) unregister() {
	close(c.msg)
	c.room.unregister <- c
	c.conn.Close()
}

func (c *client) read() {
	c.conn.SetReadLimit(maxMessageSize)
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
			c.cancel()
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
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				c.cancel()
			}

			w.Write(msg)

			if qued := len(c.msg); qued > 0 {
				for i := 0; i < qued; i++ {
					w.Write(newLine)
					w.Write(<-c.msg)
				}
			}

			if err := w.Close(); err != nil {
				c.cancel()
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.cancel()
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
