package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"server/accounts"
	"server/models"

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
		log.Fatal(err)
	}

	authedUser, err := accounts.Login(user.Username, user.Password)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusForbidden)
	} else {
		cookie := &http.Cookie{
			Name:     "authToken",
			Value:    authedUser.Token,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
		}
		http.SetCookie(c.Writer, cookie)
		json.NewEncoder(c.Writer).Encode(authedUser)
	}
}
