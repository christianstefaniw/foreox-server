package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	accounts "server/src/apps/accounts/models"
	"server/src/apps/authentication/services"
	"server/src/errors"

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

func Login(c *gin.Context) {
	user := new(accounts.User)

	err := json.NewDecoder(c.Request.Body).Decode(user)
	if err != nil {
		log.Fatal(errors.Wrap(err, err.Error()))
	}

	authedUser, err := services.Login(user.Username, user.Password)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusForbidden)
		return
	}
	cookie := &http.Cookie{
		Name:     "authToken",
		Value:    authedUser.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}

	http.SetCookie(c.Writer, cookie)
	c.JSON(http.StatusOK, authedUser)
}

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
