package server

import (
	"fmt"
	"net/http"

	"github.com/Shambou/golang-challenge/internal/objects"
	"github.com/Shambou/golang-challenge/internal/validator"
	"github.com/gorilla/mux"
)

// GetLatestFileRate - gets the latest requested rate for quote_currency
func (h *Handler) GetLatestFileRate(w http.ResponseWriter, r *http.Request) {
	v := validator.New(mux.Vars(r))
	v.Length("quote_currency", 3)

	if !v.Valid() {
		fmt.Println(v.Errors)
		jsonResponse(w, http.StatusBadRequest, "Invalid request", nil, v.Errors)
		return
	}

	currencyRate, err := h.File.GetLastRate(r.Context(), v.Get("quote_currency"))
	if err != nil {
		jsonResponse(w, http.StatusOK, err.Error(), nil, nil)
		return
	}

	message := fmt.Sprintf("Last rate for stored for %s%s", currencyRate.QuoteCurrency, currencyRate.BaseCurrency)

	jsonResponse(w, http.StatusOK, message, objects.BaseRateResponse{
		Date:          currencyRate.Date.Format("2006-02-01"),
		BaseCurrency:  currencyRate.BaseCurrency,
		QuoteCurrency: currencyRate.QuoteCurrency,
		Rate:          currencyRate.Rate.StringFixedBank(4),
	}, nil)
}
