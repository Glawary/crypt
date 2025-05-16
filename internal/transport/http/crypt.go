package http

import (
	"encoding/json"
	"net/http"

	"github.com/Glawary/crypt/internal/usecase/model"
)

func (rec *Server) ListCrypto(w http.ResponseWriter, r *http.Request) {
	var filter model.Filter
	err := json.NewDecoder(r.Body).Decode(&filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := rec.cryptService.ListCryptoCurrency(r.Context(), &filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}
}
