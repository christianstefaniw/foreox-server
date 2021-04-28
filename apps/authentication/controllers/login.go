package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	accounts "server/apps/accounts/models"
	"server/apps/authentication/services"
	"server/errors"

	"github.com/gin-gonic/gin"
)

// Logs in user
func Login(c *gin.Context) {
	c.Header("Content-Type", "application/json")
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
	json.NewEncoder(c.Writer).Encode(authedUser)

}