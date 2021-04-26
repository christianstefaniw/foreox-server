package messaging

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Room struct {
	id         primitive.ObjectID
	clients    map[*client]bool
	broadcast  chan []byte
	register   chan *client
	unregister chan *client
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewRoom() *Room {
	ctx, cancel := context.WithCancel(context.Background())
	return &Room{
		id:         primitive.NewObjectID(),
		clients:    make(map[*client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *client),
		unregister: make(chan *client),
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (r *Room) Serve() {
	for {
		select {
		case msg := <-r.broadcast:
			for c := range r.clients {
				c.msg <- msg
			}
		case client := <-r.register:
			r.clients[client] = true

		case client := <-r.unregister:
			delete(r.clients, client)
		}
	}
}
