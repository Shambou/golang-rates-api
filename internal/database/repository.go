package database

import (
	"context"
	"time"

	"github.com/Shambou/golang-challenge/internal/models"
)

// DatabaseRepo - contract for our DB calls
type DatabaseRepo interface {
	CreateRate(ctx context.Context, rate *models.CurrencyRate) error
	GetLastRate(ctx context.Context, quoteCurrency string) (models.CurrencyRate, error)
	GetRatesInRange(ctx context.Context, quoteCurrency string, fromDate time.Time, toDate time.Time) ([]models.CurrencyRate, error)
	GetAllRatesOnDate(ctx context.Context, date time.Time) ([]models.CurrencyRate, error)
	CheckRateQuoteOnDateExists(ctx context.Context, quoteCurrency string, date time.Time) bool
	TableSeeded(ctx context.Context) bool
	Ping(ctx context.Context) error
}
