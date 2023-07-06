package server

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func NewServer(address string, router *chi.Mux) *http.Server {
	return &http.Server{Addr: address, Handler: router}
}
