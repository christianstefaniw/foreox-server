package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AccountInfo(c *gin.Context) {
	// TODO handle error
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, user)
}
