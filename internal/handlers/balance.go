package handlers

import "net/http"

func (h *Handler) GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	login, status := h.getUsernameFromToken(r)
	if status != http.StatusOK {
		w.WriteHeader(status)
		return
	}

	userBalance, err := h.db.GetBalanceInfo(login)
	if err != nil {
		h.log.Errorf("error while getting user balance from db: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(userBalance)
	if err != nil {
		h.log.Errorf("error write user balance is: %s", err)
	}
}
