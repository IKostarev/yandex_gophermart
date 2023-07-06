package handlers

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
)

func (h *Handler) BasicAuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			h.log.Errorf("token is empty")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tkn, err := h.extractJwtToken(r)
		if err != nil {
			if errors.Is(err, jwt.ErrSignatureInvalid) ||
				errors.Is(err, jwt.ErrTokenExpired) ||
				errors.Is(err, ErrTokenIsEmpty) ||
				errors.Is(err, ErrNoToken) {
				h.log.Errorf(err.Error())
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			h.log.Errorf(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !tkn.Valid {
			h.log.Errorf("invalid token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.Header().Add("Authorization", tokenHeader)
		next.ServeHTTP(w, r)
	})
}
