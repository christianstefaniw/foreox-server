package messaging

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type room struct {
	id         primitive.ObjectID
	clients    map[*client]bool
	broadcast  chan []byte
	register   chan *client
	unregister chan *client
}

func NewRoom() *room {
	return &room{
		id:         primitive.NewObjectID(),
		clients:    make(map[*client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *client),
		unregister: make(chan *client),
	}
}

func (r *room) Run() {
	r.serve()
}

func (r *room) serve() {
	for {
		select {
		case msg := <-r.broadcast:
			for c := range r.clients {
				// TODO this is dangerous because if the client's msg chan is blocking messages will not be propagated
				c.msg <- msg
			}
		case client := <-r.register:
			r.clients[client] = true

		case client := <-r.unregister:
			delete(r.clients, client)
		}
	}
}
