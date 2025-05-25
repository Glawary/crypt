package http

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Glawary/crypt/internal/usecase/model"
)

// ListCrypto Список криптовалют
// @Tags Crypt
// @Summary Получение информации по криптовалютам
// @ID list-crypto
// @Accept json
// @Produce json
// @Param cryptoexchange_name query string false "Название биржы"
// @Param price_from query float64 false "Цена снизу"
// @Param price_to query float64 false "Цена сверху"
// @Param find_brush query bool false "Ёршик"
// @Success 200 {object} model.Crypto "данные по криптовалюте"
// @Router /api/v1/list [get]
func (rec *Server) ListCrypto(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	filter, err := getFilter(r.URL.Query(), w)

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

func getFilter(query url.Values, w http.ResponseWriter) (model.Filter, error) {
	var err error
	filter := model.Filter{}

	exchangeName := query.Get("cryptoexchange_name")
	if exchangeName != "" {
		filter.CryptoExchangeName = exchangeName
	}

	var priceFromFloat float64
	priceFrom := query.Get("price_from")
	if priceFrom != "" {
		priceFromFloat, err = strconv.ParseFloat(priceFrom, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		filter.PriceFrom = priceFromFloat
	}

	var priceToFloat float64
	priceTo := query.Get("price_to")
	if priceTo != "" {
		priceToFloat, err = strconv.ParseFloat(priceTo, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		filter.PriceTo = priceToFloat
	}

	findBrush := query.Get("find_brush")
	if findBrush != "" {
		if findBrush == "true" {
			filter.FindBrush = true
		} else if findBrush == "false" {
			filter.FindBrush = false
		}
	}
	return filter, nil
}
