package http

import (
	"encoding/json"
	"net/http"

	"github.com/Glawary/crypt/internal/usecase/model"
)

// ListCrypto Список криптовалют
// @Tags Crypt
// @Summary Получение информации по криптовалютам
// @ID list-crypto
// @Accept json
// @Produce json
// @Param cryptoexchange_name query string false "Название биржы"
// @Success 200 {object} model.Crypto "данные по криптовалюте"
// @Router /api/v1/list [get]
func (rec *Server) ListCrypto(w http.ResponseWriter, r *http.Request) {
	filter := model.Filter{}

	exchangeName := r.URL.Query().Get("cryptoexchange_name")
	if exchangeName != "" {
		filter.CryptexchangeName = exchangeName
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
