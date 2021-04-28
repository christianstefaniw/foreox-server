package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"server/accounts"
	"server/errors"
	"server/models"
	"time"

	"github.com/gin-gonic/gin"
)

// Registers user
func RegisterHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	user := new(models.User)
	err := json.NewDecoder(c.Request.Body).Decode(user)

	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	err = accounts.Register(user)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusForbidden)
		return
	}
	json.NewEncoder(c.Writer).Encode(user)

}

// Logs in user
func LoginHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	user := new(models.User)

	err := json.NewDecoder(c.Request.Body).Decode(user)
	if err != nil {
		log.Fatal(errors.Wrap(err, err.Error()))
	}

	authedUser, err := accounts.Login(user.Username, user.Password)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusForbidden)
		return
	}
	cookie := &http.Cookie{
		Name:  "authToken",
		Value: authedUser.Token,
		Path:  "/",
		// HttpOnly: true,
		HttpOnly: false,
		Secure:   true,
		Expires:  time.Now().Add(time.Hour * 168),
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(c.Writer, cookie)
	json.NewEncoder(c.Writer).Encode(authedUser)

}
