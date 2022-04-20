package database

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/Shambou/golang-challenge/internal/models"
	"github.com/shopspring/decimal"
)

type CSV [][]string

// CreateRate - creates new rate in db
func (f *File) CreateRate(ctx context.Context, rate *models.CurrencyRate) error {
	return nil
}

// GetLastRate - gets last rate available for
func (f *File) GetLastRate(ctx context.Context, quoteCurrency string) (models.CurrencyRate, error) {
	symbol := strings.ToTitle(quoteCurrency) + f.BaseCurrency
	ratePath := f.FxPath + symbol + f.Ext

	csvFile, err := os.Open(ratePath)
	if err != nil {
		log.Println(err)
		return models.CurrencyRate{}, err
	}
	defer csvFile.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(csvFile)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Println(err)
		return models.CurrencyRate{}, err
	}

	data = data[1:]

	sort.Sort(CSV(data))

	date, err := time.Parse("2006-01-02", data[0][0])
	if err != nil {
		return models.CurrencyRate{}, errors.New("invalid date")
	}
	rate, err := decimal.NewFromString(data[0][1])
	if err != nil {
		return models.CurrencyRate{}, errors.New("invalid rate")
	}

	currencyRate := models.CurrencyRate{
		BaseCurrency:  f.BaseCurrency,
		QuoteCurrency: strings.ToTitle(quoteCurrency),
		Date:          date,
		Rate:          rate,
	}

	fmt.Println("rate object ", currencyRate)

	return currencyRate, nil
}

func (f *File) GetRatesInRange(ctx context.Context, quoteCurrency string, fromDate time.Time, toDate time.Time) ([]models.CurrencyRate, error) {
	var rates []models.CurrencyRate
	return rates, nil
}

// CheckRateQuoteOnDateExists - Checks if rate exists in db
func (f *File) CheckRateQuoteOnDateExists(ctx context.Context, quoteCurrency string, date time.Time) bool {
	return false
}

// GetAllRatesOnDate - gets all available rates on date
func (f *File) GetAllRatesOnDate(ctx context.Context, date time.Time) ([]models.CurrencyRate, error) {
	var rates []models.CurrencyRate
	return rates, nil
}

// TableSeeded - checks if db table is already seeded
func (f *File) TableSeeded(ctx context.Context) bool {
	return true
}

// Other functions required for sort.Sort.
func (data CSV) Len() int {
	return len(data)
}
func (data CSV) Swap(i, j int) {
	data[i], data[j] = data[j], data[i]
}

func (data CSV) Less(i, j int) bool {
	dateColumnIndex := 0
	date1 := data[i][dateColumnIndex]
	date2 := data[j][dateColumnIndex]
	timeT1, _ := time.Parse("2006-01-02", date1)
	timeT2, _ := time.Parse("2006-01-02", date2)

	return timeT1.After(timeT2)
}
