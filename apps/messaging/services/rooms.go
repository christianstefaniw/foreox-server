package services

import (
	"context"
	"fmt"
	"os"
	"server/constants"
	"server/database"
	"server/errors"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Room struct {
	Id         primitive.ObjectID `bson:"_id" json:"id"`
	clients    map[*client]bool
	broadcast  chan []byte
	register   chan *client
	unregister chan *client
	Name       string `json:"name"`
	ctx        context.Context
	cancel     context.CancelFunc
}

var activeRooms sync.Map

func GetRoom(id string) (*Room, bool) {
	rm, ok := activeRooms.Load(id)
	return rm.(*Room), ok
}

func NewRoom(name string) *Room {
	ctx, cancel := context.WithCancel(context.Background())
	rm := &Room{
		Id:         primitive.NewObjectID(),
		Name:       name,
		clients:    make(map[*client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *client),
		unregister: make(chan *client),
		ctx:        ctx,
		cancel:     cancel,
	}

	err := rm.save()
	if err != nil {
		fmt.Fprint(os.Stderr, errors.Wrap(err, err.Error()))
		return nil
	}

	return rm
}

func (r *Room) roomFromDatabase() {
	ctx, cancel := context.WithCancel(context.Background())
	r.clients = make(map[*client]bool)
	r.broadcast = make(chan []byte)
	r.register = make(chan *client)
	r.unregister = make(chan *client)
	r.ctx = ctx
	r.cancel = cancel
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

func (r *Room) save() error {
	r.saveToActiveRooms()
	return r.saveToDb()
}

func (r *Room) saveToActiveRooms() {
	activeRooms.LoadOrStore(r.Id.Hex(), r)
}

func (r *Room) saveToDb() error {
	_, err := database.Database.InsertOne(context.Background(), constants.ROOMS_COLL, r)
	return err
}

func (r *Room) StartFromDb() {
	r.roomFromDatabase()
	r.saveToActiveRooms()
	go r.Serve()
}
