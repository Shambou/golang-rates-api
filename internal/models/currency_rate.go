package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type CurrencyRate struct {
	ID            int             `json:"id"`
	BaseCurrency  string          `json:"base_currency"`
	QuoteCurrency string          `json:"quote_currency"`
	Rate          decimal.Decimal `json:"rate"`
	Date          time.Time       `json:"date"`
}
