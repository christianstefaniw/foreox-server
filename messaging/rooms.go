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
	name       string
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewRoom() *Room {
	ctx, cancel := context.WithCancel(context.Background())
	return &Room{
		id:         primitive.NewObjectID(),
		name:       "new room",
		clients:    make(map[*client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *client),
		unregister: make(chan *client),
		ctx:        ctx,
		cancel:     cancel,
	}
}

func GetRoom(id primitive.ObjectID) *Room {
	return nil
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
