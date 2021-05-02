package helpers

import "context"

type Room interface {
	Serve()
	Save()
	Unregister()
	Broadcast(msg []byte)
	GetCtx() context.Context
}

type Client interface {
	Unregister()
	Read()
	Write()
	DoWork()
	ReceiveMessage(msg []byte)
}
