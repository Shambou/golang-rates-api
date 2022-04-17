package test

import (
	"fmt"
	"testing"

	"github.com/Shambou/golang-challenge/internal/server"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

const (
	BaseUrl = "http://localhost:8080/api/v1/rates"
)

func TestGetLatestRate(t *testing.T) {
	client := resty.New()
	jsonResp := &server.JsonResponse{}

	t.Run("test get rate:valid", func(t *testing.T) {
		resp, err := client.R().
			SetQueryParam("quote_currency", "chf").
			SetResult(jsonResp).
			Get(BaseUrl + "/latest")

		assert.NoError(t, err)

		assert.Equal(t, 200, resp.StatusCode())
		assert.Equal(t, fmt.Sprintf("Last rate for stored for CHFUSD"), jsonResp.Message)
	})

	t.Run("test get rate:currency not found", func(t *testing.T) {
		resp, err := client.R().
			SetQueryParam("quote_currency", "asd").
			SetResult(jsonResp).
			Get(BaseUrl + "/latest")

		assert.NoError(t, err)

		assert.Equal(t, 200, resp.StatusCode())
		assert.Equal(t, "could not get rate for ASD", jsonResp.Message)
	})

	t.Run("test get rate:invalid currency", func(t *testing.T) {
		resp, err := client.R().
			SetQueryParam("quote_currency", "").
			Get(BaseUrl + "/latest")

		assert.NoError(t, err)

		assert.Equal(t, 400, resp.StatusCode())
	})
}

func TestGetRatesInRange(t *testing.T) {
	client := resty.New()
	jsonResp := &server.JsonResponse{}
	t.Run("test get range rates:valid", func(t *testing.T) {
		resp, err := client.R().
			SetQueryParam("quote_currency", "chf").
			SetQueryParam("from", "2016-01-30").
			SetQueryParam("to", "2016-02-03").
			SetResult(jsonResp).
			Get(BaseUrl + "/range")

		assert.NoError(t, err)

		assert.Equal(t, 200, resp.StatusCode())
		assert.Equal(t, fmt.Sprintf("Rates CHFUSD in range 2016-01-30:2016-02-03"), jsonResp.Message)
	})

	t.Run("test get range rates:invalid currency", func(t *testing.T) {
		resp, err := client.R().
			SetQueryParam("quote_currency", "AS").
			SetQueryParam("from", "2016-01-30").
			SetQueryParam("to", "2016-02-03").
			Get(BaseUrl + "/range")

		assert.NoError(t, err)

		assert.Equal(t, 400, resp.StatusCode())
	})

	t.Run("test get range rates:invalid from date", func(t *testing.T) {
		resp, err := client.R().
			SetQueryParam("quote_currency", "CHF").
			SetQueryParam("from", "2016-30-30").
			SetQueryParam("to", "2016-02-03").
			Get(BaseUrl + "/range")

		assert.NoError(t, err)

		assert.Equal(t, 400, resp.StatusCode())
	})

	t.Run("test get range rates:invalid to date", func(t *testing.T) {
		resp, err := client.R().
			SetQueryParam("quote_currency", "CHF").
			SetQueryParam("from", "2016-01-30").
			SetQueryParam("to", "2016-30-03").
			Get(BaseUrl + "/range")

		assert.NoError(t, err)

		assert.Equal(t, 400, resp.StatusCode())
	})
}

func TestGetTimeseriesData(t *testing.T) {
	client := resty.New()
	jsonResp := &server.JsonResponse{}
	t.Run("test get timeseries rates:valid", func(t *testing.T) {
		resp, err := client.R().
			SetQueryParam("date", "2016-04-13").
			SetResult(jsonResp).
			Get(BaseUrl + "/timeseries")

		assert.NoError(t, err)

		assert.Equal(t, 200, resp.StatusCode())
		assert.Equal(t, fmt.Sprintf("All available currency rates on 2016-04-13"), jsonResp.Message)
	})

	t.Run("test get timeseries rates:invalid date", func(t *testing.T) {
		resp, err := client.R().
			SetQueryParam("date", "2016-30-30").
			Get(BaseUrl + "/timeseries")

		assert.NoError(t, err)

		assert.Equal(t, 400, resp.StatusCode())
	})
}

func TestStoreRate(t *testing.T) {
	client := resty.New()
	jsonResp := &server.JsonResponse{}

	t.Run("test store rate:rate already exists", func(t *testing.T) {
		resp, err := client.R().
			SetBody(`{"date": "2016-01-29","rate": "1.022600"}`).
			SetResult(jsonResp).
			Post(BaseUrl + "/chf")

		assert.NoError(t, err)

		assert.Equal(t, 422, resp.StatusCode())
	})

	t.Run("test store rate:invalid date", func(t *testing.T) {
		resp, err := client.R().
			SetBody(`{"date": "2022-05-29","rate": "1.022600"}`).
			Post(BaseUrl + "/chf")

		assert.NoError(t, err)

		assert.Equal(t, 422, resp.StatusCode())
	})

	t.Run("test store rate:invalid rate", func(t *testing.T) {
		resp, err := client.R().
			SetBody(`{"date": "2022-01-29","rate": "-10"}`).
			Post(BaseUrl + "/chf")

		assert.NoError(t, err)

		assert.Equal(t, 400, resp.StatusCode())
	})
}
