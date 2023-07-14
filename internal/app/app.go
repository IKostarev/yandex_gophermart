package app

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"yandex_gophermart/internal/handlers"
	storage "yandex_gophermart/internal/storage/db"
)

func NewApp(dbManager *storage.Manager, log *zap.SugaredLogger) *chi.Mux {
	handler := handlers.NewHandlers(dbManager, log)

	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Post("/api/user/register", handler.RegisterHandler)
		r.Post("/api/user/login", handler.LoginHandler)
	})

	r.Group(func(r chi.Router) {
		r.Use(handler.BasicAuthHandler)
		r.Post("/api/user/orders", handler.LoadOrderHandler)
		r.Post("/api/user/balance/withdraw", handler.WithdrawHandler)
		r.Get("/api/user/orders", handler.GetOrdersHandler)
		r.Get("/api/user/withdrawals", handler.GetWithdrawalsHandler)
		r.Get("/api/user/balance", handler.GetBalanceHandler)
	})

	return r
}
