package seeds

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	database "github.com/Shambou/golang-challenge/internal/database/postgres"
	"github.com/Shambou/golang-challenge/internal/models"
	"github.com/shopspring/decimal"
)

// Seed type
type Seed struct {
	DB *database.Database
}

func New(db *database.Database) *Seed {
	s := &Seed{
		DB: db,
	}

	return s
}

func (s *Seed) Execute() {
	if s.DB.TableSeeded(context.Background()) {
		log.Println("Table already seeded")
		return
	}
	ext := ".csv"
	err := filepath.WalkDir("fxdata", func(str string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			s.seed(str)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
}

// seed - Parse CSV and seed the db
// code taken from https://github.com/jasonbronson/example_Go_file_import/blob/master/part3/main.go
func (s *Seed) seed(path string) {
	fmt.Println("found seed in path ", path)

	symbol := strings.Split(path, "/")
	// get base and quote currency from symbol slice
	quoteCurrency := symbol[1][0:3]
	baseCurrency := symbol[1][3:6]

	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(f)
	var line string
	var lineCount int64
	var rates []models.CurrencyRate
	for {
		line, err = reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		if lineCount == 0 {
			lineCount++
			continue
		}
		lineCount++
		r := csv.NewReader(strings.NewReader(line))
		r.Comma = ','
		records, err := r.ReadAll()
		if err != nil {
			panic(err)
		}

		var rate models.CurrencyRate
		rate.BaseCurrency = baseCurrency
		rate.QuoteCurrency = quoteCurrency

		for _, row := range records {
			rate.Date, err = time.Parse("2006-01-02", row[0])
			if err != nil {
				continue
			}
			rate.Rate, err = decimal.NewFromString(row[1])
			if err != nil {
				continue
			}
			rates = append(rates, rate)
		}
	}

	if err = s.DB.BulkInsert(rates); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Finished")
}
