package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"yandex_gophermart/internal/models"
	"yandex_gophermart/internal/storage/db"
)

func (h *Handler) WithdrawHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	var withdrawInfo *models.WithdrawInfo
	var buf bytes.Buffer

	if _, err := buf.ReadFrom(r.Body); err != nil {
		h.log.Errorf("error while reading request body: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &withdrawInfo); err != nil {
		h.log.Errorf("error while unmarshalling request body: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !h.checkOrder(withdrawInfo.OrderID) {
		h.log.Error("invalid order format")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	login, status := h.getUsernameFromToken(r)
	if status != http.StatusOK {
		w.WriteHeader(status)
		return
	}

	if err := h.db.Withdraw(login, withdrawInfo.OrderID, withdrawInfo.Amount); err != nil {
		if errors.Is(err, storage.ErrInsufficientBalance) {
			w.WriteHeader(http.StatusPaymentRequired)
			return
		}

		h.log.Errorf("error while trying to withdraw %f from user %q: %s", withdrawInfo.Amount, login, err.Error())

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	h.log.Infof("withdrawn %f from user %q for order %q", withdrawInfo.Amount, login, withdrawInfo.OrderID)
}
