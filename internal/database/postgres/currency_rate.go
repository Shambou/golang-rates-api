package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Shambou/golang-challenge/internal/models"
)

const BaseCurrency = "USD"

// CreateRate - creates new rate in db
func (d *Database) CreateRate(ctx context.Context, rate *models.CurrencyRate) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `insert into currency_rates
		(base_currency, quote_currency, rate, date) VALUES
		($1, $2, $3, $4) returning id`

	rate.BaseCurrency = BaseCurrency

	result, err := d.Client.ExecContext(
		ctx,
		query,
		rate.BaseCurrency,
		rate.QuoteCurrency,
		rate.Rate,
		rate.Date,
	)
	if err != nil {
		log.Println(err)
		return err
	}
	id, _ := result.LastInsertId()
	rate.ID = int(id)

	return nil
}

// GetLastRate - gets last rate available for
func (d *Database) GetLastRate(ctx context.Context, quoteCurrency string) (models.CurrencyRate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	var rate = models.CurrencyRate{}
	quoteCurrency = strings.ToTitle(quoteCurrency)

	row := d.Client.QueryRowContext(
		ctx,
		"select date, base_currency, quote_currency, rate from currency_rates where quote_currency = $1 order by date desc limit 1",
		quoteCurrency,
	)
	err := row.Scan(&rate.Date, &rate.BaseCurrency, &rate.QuoteCurrency, &rate.Rate)
	if err != nil {
		return rate, errors.New(fmt.Sprintf("could not get rate for %s", quoteCurrency))
	}

	return rate, nil
}

func (d *Database) GetRatesInRange(ctx context.Context, quoteCurrency string, fromDate time.Time, toDate time.Time) ([]models.CurrencyRate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var rates []models.CurrencyRate

	from := fromDate.Format("2006-01-02")
	to := toDate.Format("2006-01-02")
	quoteCurrency = strings.ToTitle(quoteCurrency)

	query := `select date, base_currency, quote_currency, rate from currency_rates where quote_currency = $1 and date between $2 and $3 order by date asc`

	rows, err := d.Client.QueryContext(
		ctx,
		query,
		quoteCurrency,
		from,
		to,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rate models.CurrencyRate
		err := rows.Scan(
			&rate.Date,
			&rate.BaseCurrency,
			&rate.QuoteCurrency,
			&rate.Rate,
		)
		if err != nil {
			return nil, err
		}

		rates = append(rates, rate)
	}

	return rates, nil
}

// CheckRateQuoteOnDateExists - Checks if rate exists in db
func (d *Database) CheckRateQuoteOnDateExists(ctx context.Context, quoteCurrency string, date time.Time) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	searchDate := date.Format("2006-01-02")
	quoteCurrency = strings.ToTitle(quoteCurrency)

	row := d.Client.QueryRowContext(
		ctx,
		"select id from currency_rates where quote_currency = $1 and date = $2 limit 1",
		quoteCurrency,
		searchDate,
	)

	var id interface{}

	if err := row.Scan(&id); err != nil {
		return false
	}

	return true
}

// GetAllRatesOnDate - gets all available rates on date
func (d *Database) GetAllRatesOnDate(ctx context.Context, date time.Time) ([]models.CurrencyRate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var rates []models.CurrencyRate

	searchDate := date.Format("2006-01-02")

	query := `select date, base_currency, quote_currency, rate 
		from currency_rates  
		where date = $1   
		order by quote_currency asc`

	rows, err := d.Client.QueryContext(
		ctx,
		query,
		searchDate,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rate models.CurrencyRate
		err := rows.Scan(&rate.Date, &rate.BaseCurrency, &rate.QuoteCurrency, &rate.Rate)
		if err != nil {
			return nil, err
		}
		rates = append(rates, rate)
	}

	return rates, nil
}

// TableSeeded - checks if db table is already seeded
func (d *Database) TableSeeded(ctx context.Context) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := d.Client.QueryRowContext(
		ctx,
		"select 1 from currency_rates limit 1",
	)
	var result interface{}
	if err := row.Scan(&result); err != nil {
		fmt.Println(err)
		fmt.Println("opet")
		return false
	}
	return true
}

// BulkInsert - used only to bulk insert data from csv files
// code taken from https://stackoverflow.com/questions/60026385/bulk-insert-with-sqlx
func (d *Database) BulkInsert(rates []models.CurrencyRate) error {
	queryInsert := `INSERT INTO currency_rates (base_currency, quote_currency, rate, date) VALUES `
	insertparams := []interface{}{}
	for i, rate := range rates {
		p1 := i * 4 // starting position for insert params
		queryInsert += fmt.Sprintf("($%d,$%d,$%d,$%d),", p1+1, p1+2, p1+3, p1+4)
		insertparams = append(insertparams, rate.BaseCurrency, rate.QuoteCurrency, rate.Rate, rate.Date)
	}
	queryInsert = queryInsert[:len(queryInsert)-1] // remove trailing ","
	d.Client.MustExec(queryInsert, insertparams...)

	return nil
}
