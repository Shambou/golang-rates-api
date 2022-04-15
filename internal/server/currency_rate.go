package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	database "github.com/Shambou/golang-challenge/internal/database/postgres"
	"github.com/Shambou/golang-challenge/internal/models"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
)

type JsonDateRate struct {
	Date string `json:"date"`
	Rate string `json:"rate"`
}

type JsonQuoteRate struct {
	QuoteCurrency string `json:"quote_currency"`
	Rate          string `json:"rate"`
}

type PostRateRequest struct {
	Date string
	Rate decimal.Decimal
}

// GetLatestRate - gets the latest requested rate for quote_currency
func (h *Handler) GetLatestRate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if len(vars["quote_currency"]) != 3 {
		jsonResponse(w, http.StatusBadRequest, "Quote currency not valid", nil)
		return
	}

	currencyRate, err := h.DB.GetLastRate(r.Context(), vars["quote_currency"])
	if err != nil {
		fmt.Println(err)
		jsonResponse(w, http.StatusOK, err.Error(), nil)
		return
	}

	data := make(map[string]interface{})
	data["date"] = currencyRate.Date.Format("2006-02-01")
	data["base_currency"] = currencyRate.BaseCurrency
	data["quote_currency"] = currencyRate.QuoteCurrency
	data["rate"] = currencyRate.Rate.StringFixedBank(4)
	message := fmt.Sprintf("Last rate for stored for %s%s", currencyRate.QuoteCurrency, currencyRate.BaseCurrency)

	jsonResponse(w, http.StatusOK, message, data)
}

// GetRatesInRange - gets the rates between two dates
func (h *Handler) GetRatesInRange(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if len(vars["quote_currency"]) != 3 {
		jsonResponse(w, http.StatusBadRequest, "Quote currency not valid", nil)
		return
	}
	fromDate, err := time.Parse("2006-01-02", vars["from"])
	if err != nil {
		fmt.Println(err)
		jsonResponse(w, http.StatusBadRequest, "From date is not valid", nil)
		return
	}
	toDate, err := time.Parse("2006-01-02", vars["to"])
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, "To date is not valid", nil)
		return
	}

	rates, err := h.DB.GetRatesInRange(r.Context(), vars["quote_currency"], fromDate, toDate)
	if err != nil {
		fmt.Println(err)
		jsonResponse(w, http.StatusOK, err.Error(), nil)
		return
	}

	var rangeRates []JsonDateRate

	for _, r := range rates {
		rate := JsonDateRate{}
		rate.Rate = r.Rate.StringFixedBank(4)
		rate.Date = r.Date.Format("2006-01-02")

		rangeRates = append(rangeRates, rate)
	}

	data := make(map[string]interface{})
	data["base_currency"] = database.BaseCurrency
	data["quote_currency"] = strings.ToTitle(vars["quote_currency"])
	data["rates"] = rangeRates

	message := fmt.Sprintf(
		"Rates %s%s in range %s:%s", data["quote_currency"],
		database.BaseCurrency,
		fromDate.Format("2006-01-02"),
		toDate.Format("2006-01-02"),
	)

	jsonResponse(w, http.StatusOK, message, data)
}

// GetTimeseriesData - gets the all available rates on date
func (h *Handler) GetTimeseriesData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	date, err := time.Parse("2006-01-02", vars["date"])
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, "Date is not valid", nil)
		return
	}

	rates, err := h.DB.GetAllRatesOnDate(r.Context(), date)
	if err != nil {
		fmt.Println(err)
		jsonResponse(w, http.StatusOK, err.Error(), nil)
		return
	}

	var quoteRates []JsonQuoteRate

	for _, r := range rates {
		rate := JsonQuoteRate{}
		rate.Rate = r.Rate.StringFixedBank(4)
		rate.QuoteCurrency = r.QuoteCurrency

		quoteRates = append(quoteRates, rate)
	}

	data := make(map[string]interface{})
	data["base_currency"] = database.BaseCurrency
	data["date"] = date
	data["rates"] = quoteRates

	message := fmt.Sprintf("All available currency rates on %s", date.Format("2006-01-02"))

	jsonResponse(w, http.StatusOK, message, data)
}

// StoreRate - stores new rate
func (h *Handler) StoreRate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var postRateReq PostRateRequest

	if err := json.NewDecoder(r.Body).Decode(&postRateReq); err != nil {
		jsonResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	date, err := time.Parse("2006-01-02", postRateReq.Date)
	if err != nil || date.After(time.Now()) {
		jsonResponse(w, http.StatusBadRequest, "Date is not valid", nil)
		return
	}
	rate := postRateReq.Rate
	if rate.LessThanOrEqual(decimal.Zero) {
		jsonResponse(w, http.StatusBadRequest, "Rate is not valid", nil)
		return
	}

	if h.DB.CheckRateQuoteOnDateExists(r.Context(), vars["currency"], date) {
		jsonResponse(w, http.StatusUnprocessableEntity, "Rate for this currency and date already exists", nil)
		return
	}

	var currencyRate = models.CurrencyRate{}
	currencyRate.QuoteCurrency = strings.ToTitle(vars["currency"])
	currencyRate.Date = date
	currencyRate.Rate = rate

	err = h.DB.CreateRate(r.Context(), &currencyRate)
	if err != nil {
		fmt.Println(err)
		jsonResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	data := make(map[string]interface{})
	data["base_currency"] = currencyRate.BaseCurrency
	data["quote_currency"] = currencyRate.QuoteCurrency
	data["rate"] = currencyRate.Rate
	data["date"] = currencyRate.Date.Format("2006-01-02")

	jsonResponse(w, http.StatusCreated, "Stored new rate", data)
}
