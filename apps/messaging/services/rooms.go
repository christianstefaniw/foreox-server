package services

import (
	"context"
	"sync"

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

var activeRooms sync.Map

func GetRoom(id string) (*Room, bool) {
	rm, ok := activeRooms.Load(id)
	return rm.(*Room), ok
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

func (r *Room) Serve(rmId chan<- string) {
	rmId <- r.id.Hex()
	activeRooms.LoadOrStore(r.id.Hex(), r)
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
