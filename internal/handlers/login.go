package handlers

import (
	"fmt"
	"net/http"
	"time"
)

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	user, success := h.parseInputUser(r)
	if !success {
		h.log.Errorf("parse input user is error: %s", user)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.db.Login(user.Login, user.Password); err != nil {
		h.log.Errorf("error while login user: %s", err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(time.Hour)

	token, err := createToken(user.Login, expirationTime)
	if err != nil {
		h.log.Errorf("error while create token for user: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Authorization", fmt.Sprintf("Bearer %s", token))

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: expirationTime,
	})

	h.log.Info(fmt.Sprintf("user %q is successfully authorized", user.Login))
}
