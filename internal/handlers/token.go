package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
	"time"
	"yandex_gophermart/internal/models"
)

var jwtKey = []byte("my_secret_key")

func (h *Handler) extractJwtToken(r *http.Request) (*jwt.Token, error) {
	tokenHeader := r.Header.Get("Authorization")

	if tokenHeader == "" {
		h.log.Errorf("token is empty")
		return nil, ErrTokenIsEmpty
	}

	splitted := strings.Split(tokenHeader, " ")
	if len(splitted) != 2 {
		h.log.Errorf("no token")
		return nil, ErrNoToken
	}

	tknStr := splitted[1]

	claims := &models.Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parse with claims is: %s", err)
	}

	return tkn, err
}

func (h *Handler) parseInputUser(r *http.Request) (*models.User, bool) {
	var userFromRequest *models.User
	var buf bytes.Buffer

	if _, err := buf.ReadFrom(r.Body); err != nil {
		h.log.Errorf("error while reading request body: %s", err.Error())
		return nil, false
	}

	if err := json.Unmarshal(buf.Bytes(), &userFromRequest); err != nil {
		h.log.Errorf("error while unmarshalling request body: %s", err.Error())
		return nil, false
	}

	if userFromRequest.Login == "" || userFromRequest.Password == "" {
		h.log.Errorf("login or password is empty")
		return nil, false
	}

	return userFromRequest, true
}

func (h *Handler) getUsernameFromToken(r *http.Request) (string, int) {
	var data bytes.Buffer

	if _, err := data.ReadFrom(r.Body); err != nil {
		h.log.Errorf("error while reading request body: %s", err.Error())
		return "", http.StatusBadRequest
	}

	tkn, err := h.extractJwtToken(r)
	if err != nil {
		h.log.Errorf("error while extracting token: %s", err.Error())
		return "", http.StatusInternalServerError
	}

	claims, ok := tkn.Claims.(*models.Claims)
	if !ok {
		h.log.Errorf("error while getting claims")
		return "", http.StatusInternalServerError
	}

	return claims.Username, http.StatusOK
}

func createToken(userName string, expirationTime time.Time) (string, error) {
	claims := &models.Claims{
		Username: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("error signed string jwt token is: %s", err)
	}

	return tokenString, nil
}
