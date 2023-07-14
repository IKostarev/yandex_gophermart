package handlers

import (
	"errors"
	"net/http"
	"yandex_gophermart/internal/storage/db"
)

func (h *Handler) GetWithdrawalsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	login, status := h.getUsernameFromToken(r)
	if status != http.StatusOK {
		w.WriteHeader(status)
		return
	}

	userWithdrawals, err := h.db.GetWithdrawals(login)
	if err != nil {
		if errors.Is(err, storage.ErrNoData) {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h.log.Errorf("error while getting withdrawals from db: %s", err.Error())

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, err = w.Write(userWithdrawals)
	if err != nil {
		h.log.Errorf("error write user withdrawals is: %s", err)
	}
}
