package objects

type JsonDateRateResponse struct {
	Date string `json:"date"`
	Rate string `json:"rate"`
}

type JsonDateRateResponses []JsonDateRateResponse

type JsonQuoteRateResponse struct {
	QuoteCurrency string `json:"quote_currency"`
	Rate          string `json:"rate"`
}

type JsonQuoteRateResponses []JsonQuoteRateResponse

type BaseRateResponse struct {
	Date          string `json:"date"`
	BaseCurrency  string `json:"base_currency"`
	QuoteCurrency string `json:"quote_currency"`
	Rate          string `json:"rate"`
}

type RangeRatesResponse struct {
	BaseCurrency  string                `json:"base_currency"`
	QuoteCurrency string                `json:"quote_currency"`
	Rates         JsonDateRateResponses `json:"rates"`
}

type QuoteRatesResponse struct {
	BaseCurrency string                 `json:"base_currency"`
	Date         string                 `json:"date"`
	Rates        JsonQuoteRateResponses `json:"rates"`
}
