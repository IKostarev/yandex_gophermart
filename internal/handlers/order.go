package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"yandex_gophermart/internal/storage/db"
)

func (h *Handler) GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	login, status := h.getUsernameFromToken(r)
	if status != http.StatusOK {
		w.WriteHeader(status)
		return
	}

	userOrders, err := h.db.GetUserOrders(login)
	if err != nil {
		if errors.Is(err, storage.ErrNoData) {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h.log.Errorf("error while getting orders from db: %s", err.Error())

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, _ = w.Write(userOrders)
}

func (h *Handler) LoadOrderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/plain")

	var data bytes.Buffer

	if _, err := data.ReadFrom(r.Body); err != nil {
		h.log.Errorf("error while reading request body: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	login, status := h.getUsernameFromToken(r)
	if status != http.StatusOK {
		w.WriteHeader(status)
		return
	}

	order := data.String()
	if !h.checkOrder(order) {
		h.log.Error("invalid order format")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if err := h.db.LoadOrder(login, order); err != nil {
		if errors.Is(err, storage.ErrCreatedBySameUser) {
			h.log.Info(fmt.Sprintf("order %q was alredy created by the same user", order))
			w.WriteHeader(http.StatusOK)
			return
		}

		if errors.Is(err, storage.ErrCreatedDiffUser) {
			h.log.Info(fmt.Sprintf("order %q was alredy created by the other user", order))
			w.WriteHeader(http.StatusConflict)
			return
		}

		h.log.Errorf("error while loading order to db: %s", err.Error())

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) checkOrder(orderID string) bool {
	orderAsInteger, err := strconv.Atoi(orderID)
	if err != nil {
		return false
	}

	number := orderAsInteger / 10

	luhn := 0

	for i := 0; number > 0; i++ {
		c := number % 10

		if i%2 == 0 {
			c *= 2
			if c > 9 {
				c = c%10 + c/10
			}
		}

		luhn += c
		number /= 10
	}

	return (orderAsInteger%10+luhn)%10 == 0
}
