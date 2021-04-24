package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"server/accounts"
	"server/models"
)

// Registers user
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := new(models.User)
	err := json.NewDecoder(r.Body).Decode(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = accounts.Register(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	json.NewEncoder(w).Encode(user)

}

// Logs in user
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := new(models.User)

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		log.Fatal(err)
	}

	authedUser, err := accounts.Login(user.Username, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
	} else {
		c := &http.Cookie{
			Name:     "authToken",
			Value:    authedUser.Token,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
		}
		http.SetCookie(w, c)
		json.NewEncoder(w).Encode(authedUser)
	}
}
