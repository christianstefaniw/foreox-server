package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckLoggedIn(c *gin.Context) {
	// if this handler is run the user is logged in but might as well double check
	_, ok := c.Get("user")
	if !ok {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}
}
