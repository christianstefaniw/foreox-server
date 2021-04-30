package controllers

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

func AccountInfo(c *gin.Context) {
	// TODO handle error
	user, _ := c.Get("user")
	json.NewEncoder(c.Writer).Encode(user)
}
