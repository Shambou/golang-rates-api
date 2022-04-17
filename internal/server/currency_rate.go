package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	database "github.com/Shambou/golang-challenge/internal/database/postgres"
	"github.com/Shambou/golang-challenge/internal/models"
	"github.com/Shambou/golang-challenge/internal/objects"
	"github.com/Shambou/golang-challenge/internal/validator"
	"github.com/gorilla/mux"
)

// GetLatestRate - gets the latest requested rate for quote_currency
func (h *Handler) GetLatestRate(w http.ResponseWriter, r *http.Request) {
	v := validator.New(mux.Vars(r))
	v.Length("quote_currency", 3)

	if !v.Valid() {
		fmt.Println(v.Errors)
		jsonResponse(w, http.StatusBadRequest, "Invalid request", nil, v.Errors)
		return
	}

	currencyRate, err := h.DB.GetLastRate(r.Context(), v.Get("quote_currency"))
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

// GetRatesInRange - gets the rates between two dates
func (h *Handler) GetRatesInRange(w http.ResponseWriter, r *http.Request) {
	v := validator.New(mux.Vars(r))
	v.Length("quote_currency", 3)
	v.Date("from", "to")

	if !v.Valid() {
		fmt.Println(v.Errors)
		jsonResponse(w, http.StatusBadRequest, "Invalid request", nil, v.Errors)
		return
	}

	fromDate, err := time.Parse("2006-01-02", v.Get("from"))
	toDate, err := time.Parse("2006-01-02", v.Get("to"))
	quoteCurrency := strings.ToTitle(v.Get("quote_currency"))

	rates, err := h.DB.GetRatesInRange(r.Context(), quoteCurrency, fromDate, toDate)
	if err != nil {
		fmt.Println(err)
		jsonResponse(w, http.StatusOK, err.Error(), nil, nil)
		return
	}

	var rangeRates objects.JsonDateRateResponses

	for _, r := range rates {
		rate := objects.JsonDateRateResponse{}
		rate.Rate = r.Rate.StringFixedBank(4)
		rate.Date = r.Date.Format("2006-01-02")

		rangeRates = append(rangeRates, rate)
	}

	data := objects.RangeRatesResponse{
		BaseCurrency:  database.BaseCurrency,
		QuoteCurrency: quoteCurrency,
		Rates:         rangeRates,
	}

	message := fmt.Sprintf(
		"Rates %s%s in range %s:%s", data.QuoteCurrency,
		data.BaseCurrency,
		fromDate.Format("2006-01-02"),
		toDate.Format("2006-01-02"),
	)

	jsonResponse(w, http.StatusOK, message, data, nil)
}

// GetTimeseriesData - gets the all available rates on date
func (h *Handler) GetTimeseriesData(w http.ResponseWriter, r *http.Request) {
	v := validator.New(mux.Vars(r))
	v.Date("date")

	if !v.Valid() {
		fmt.Println(v.Errors)
		jsonResponse(w, http.StatusBadRequest, "Invalid request", nil, v.Errors)
		return
	}

	date, err := time.Parse("2006-01-02", v.Get("date"))
	rates, err := h.DB.GetAllRatesOnDate(r.Context(), date)

	if err != nil {
		jsonResponse(w, http.StatusOK, err.Error(), nil, nil)
		return
	}

	var quoteRates objects.JsonQuoteRateResponses

	for _, r := range rates {
		rate := objects.JsonQuoteRateResponse{}
		rate.Rate = r.Rate.StringFixedBank(4)
		rate.QuoteCurrency = r.QuoteCurrency

		quoteRates = append(quoteRates, rate)
	}

	data := objects.QuoteRatesResponse{
		BaseCurrency: database.BaseCurrency,
		Date:         date.Format("2006-01-02"),
		Rates:        quoteRates,
	}
	message := fmt.Sprintf("All available currency rates on %s", data.Date)

	jsonResponse(w, http.StatusOK, message, data, nil)
}

// StoreRate - stores new rate
func (h *Handler) StoreRate(w http.ResponseWriter, r *http.Request) {
	var postRateReq objects.PostRateRequest

	if err := json.NewDecoder(r.Body).Decode(&postRateReq); err != nil {
		jsonResponse(w, http.StatusBadRequest, err.Error(), nil, nil)
		return
	}

	vars := mux.Vars(r)
	data := make(map[string]string)
	data["currency"] = vars["currency"]
	data["rate"] = postRateReq.Rate.String()
	data["date"] = postRateReq.Date

	v := validator.New(data)
	v.Date("date")
	v.DateInFuture("date")
	v.NotEqual("currency", database.BaseCurrency)
	v.ValidRate("rate")

	if !v.Valid() {
		fmt.Println(v.Errors)
		jsonResponse(w, http.StatusBadRequest, "Invalid request", nil, v.Errors)
		return
	}

	date, err := time.Parse("2006-01-02", postRateReq.Date)

	if h.DB.CheckRateQuoteOnDateExists(r.Context(), v.Get("currency"), date) {
		jsonResponse(w, http.StatusUnprocessableEntity, "Rate for this currency and date already exists", nil, nil)
		return
	}

	var currencyRate = models.CurrencyRate{}
	currencyRate.QuoteCurrency = strings.ToTitle(vars["currency"])
	currencyRate.Date = date
	currencyRate.Rate = postRateReq.Rate

	err = h.DB.CreateRate(r.Context(), &currencyRate)
	if err != nil {
		fmt.Println(err)
		jsonResponse(w, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	jsonResponse(w, http.StatusCreated, "Stored new rate", objects.BaseRateResponse{
		BaseCurrency:  currencyRate.BaseCurrency,
		QuoteCurrency: currencyRate.QuoteCurrency,
		Rate:          currencyRate.Rate.StringFixedBank(4),
		Date:          currencyRate.Date.Format("2006-01-02"),
	}, nil)
}
