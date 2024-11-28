package server

import (
	"encoding/json"
	"net/http"
)

func HandlerTransaction(w http.ResponseWriter, r * http.Request) {
	var req TransactionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err = SendTransaction(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write([]byte("Transaction sent"))
}