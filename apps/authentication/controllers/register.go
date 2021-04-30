package controllers

import (
	"encoding/json"
	"net/http"
	accounts "server/apps/accounts/models"
	"server/apps/authentication/services"

	"github.com/gin-gonic/gin"
)

// Registers user
func Register(c *gin.Context) {
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

	c.JSON(http.StatusCreated, user)
}
