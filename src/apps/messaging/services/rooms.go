package services

import (
	"context"
	"fmt"
	"os"
	"server/src/constants"
	"server/src/database"
	"server/src/errors"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Room struct {
	Id         primitive.ObjectID `bson:"_id" json:"id"`
	clients    map[*client]bool
	broadcast  chan *message
	register   chan *client
	unregister chan *client
	Name       string     `json:"name"`
	Image      string     `json:"image"`
	Messages   []*message `json:"messages"`
	ctx        context.Context
	cancel     context.CancelFunc
}

var activeRooms sync.Map

func GetRoom(id string) (*Room, bool) {
	rm, ok := activeRooms.Load(id)
	return rm.(*Room), ok
}

func NewRoom(name string, image string) *Room {
	ctx, cancel := context.WithCancel(context.Background())
	rm := &Room{
		Id:         primitive.NewObjectID(),
		Name:       name,
		Image:      image,
		clients:    make(map[*client]bool),
		broadcast:  make(chan *message),
		register:   make(chan *client),
		unregister: make(chan *client),
		Messages:   make([]*message, 0),
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

func (r *Room) CheckClientInRoom(tkn string) (*client, bool) {
	for client := range r.clients {
		if client.Token == tkn {
			return client, true
		}
	}
	return nil, false
}

func (r *Room) roomFromDatabase() {
	ctx, cancel := context.WithCancel(context.Background())
	r.clients = make(map[*client]bool)
	r.broadcast = make(chan *message)
	r.register = make(chan *client)
	r.unregister = make(chan *client)
	r.ctx = ctx
	r.cancel = cancel
}

func (r *Room) Serve() {
	for {
		select {
		case msg := <-r.broadcast:
			// TODO error
			_ = r.saveMessage(msg)
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

func (r *Room) saveMessage(msg *message) error {
	_, err := database.Database.UpdateOne(context.Background(), constants.ROOMS_COLL, bson.M{"_id": r.Id},
		bson.M{"$push": bson.M{"messages": msg}})
	return err
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
