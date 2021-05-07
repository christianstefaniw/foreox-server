package services

import (
	accounts "server/apps/accounts/models"
	"time"
)

type message struct {
	Content string         `json:"content"`
	Time    time.Time      `json:"time"`
	Sender  *accounts.User `json:"sender"`
}

func msgStringToStruct(msg []byte, sender *accounts.User) *message {
	return &message{
		Content: string(msg),
		Time:    time.Now(),
		Sender:  sender,
	}
}
