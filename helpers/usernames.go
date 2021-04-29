package helpers

import (
	"sync"
)

var usernames sync.Map

func AddUsername(tkn, username string) {
	usernames.LoadOrStore(tkn, username)
}

func DeleteUsername(tkn string) {
	usernames.Delete(tkn)
}

func GetUsername(tkn string) string {
	username, ok := usernames.Load(tkn)
	if !ok {
		return ""
	}
	return username.(string)
}
