package controllers

import (
	"encoding/json"
	"server/helpers"

	"github.com/gin-gonic/gin"
)

type username struct {
	Username string `json:"username"`
}

func Home(c *gin.Context) {
	// TODO handle error
	tknStr, _ := c.Cookie("authToken")
	usernameStr := helpers.GetUsername(tknStr)
	json.NewEncoder(c.Writer).Encode(username{Username: usernameStr})
}
