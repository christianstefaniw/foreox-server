package controllers

import (
	"encoding/json"
	"net/http"
	"server/apps/accounts"
	"server/apps/authentication/services"

	"github.com/gin-gonic/gin"
)

// Registers user
func Register(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	user := new(accounts.User)
	err := json.NewDecoder(c.Request.Body).Decode(user)

	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	err = services.Register(user)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusForbidden)
		return
	}
	json.NewEncoder(c.Writer).Encode(user)

}
