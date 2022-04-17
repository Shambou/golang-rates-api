package objects

import "github.com/shopspring/decimal"

type PostRateRequest struct {
	Date string          `json:"date"`
	Rate decimal.Decimal `json:"rate"`
}
