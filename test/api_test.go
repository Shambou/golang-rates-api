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
		assert.Equal(t, 200, jsonResp.Status)
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
			SetResult(jsonResp).
			Get(BaseUrl + "/latest")

		assert.NoError(t, err)

		assert.Equal(t, 400, resp.StatusCode())
	})

}
