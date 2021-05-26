package settings

import (
	"context"
	rooms "server/src/apps/messaging/services"
	"server/src/constants"
	"server/src/database"

	"go.mongodb.org/mongo-driver/bson"
)

const API_PATH = "/api/"

var LoadRooms = true

func Settings() {
	if LoadRooms {
		var rooms []*rooms.Room

		ctx := context.Background()
		cursor, _ := database.Database.Find(ctx, constants.ROOMS_COLL, bson.D{{}})
		cursor.All(ctx, &rooms)
		for _, room := range rooms {
			room.StartFromDb()
		}
	}
}
