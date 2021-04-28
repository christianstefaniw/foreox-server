package controllers

import (
	"encoding/json"
	"server/foreox/settings"

	"github.com/gin-gonic/gin"
)

type username struct {
	Username string `json:"username"`
}

func Home(c *gin.Context) {
	// TODO handle error
	tknStr, _ := c.Cookie("authToken")
	usernameLoad, _ := settings.Usernames.Load(tknStr)
	json.NewEncoder(c.Writer).Encode(username{Username: usernameLoad.(string)})
}
