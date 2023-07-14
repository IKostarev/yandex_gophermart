package handlers

import (
	"go.uber.org/zap"
	"yandex_gophermart/internal/storage"
)

func NewHandlers(db storage.Storage, log *zap.SugaredLogger) *Handler {
	return &Handler{
		db:  db,
		log: log,
	}
}

type Handler struct {
	db  storage.Storage
	log *zap.SugaredLogger
}
