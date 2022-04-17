package test

import (
	"fmt"
	"testing"

	"github.com/Shambou/golang-challenge/internal/validator"
)

const BaseCurrency = "USD"

func TestValidator_Valid(t *testing.T) {
	data := make(map[string]string)
	data["currency"] = BaseCurrency
	v := validator.New(data)

	if !v.Valid() {
		t.Error("got invalid when should have been valid")
	}
}

func TestValidator_Length(t *testing.T) {
	data := make(map[string]string)
	data["currency"] = BaseCurrency

	const fakeLen = 10

	v := validator.New(data)
	v.Length("currency", fakeLen)
	if v.Valid() {
		t.Error(fmt.Sprintf("error: validator says %d length is equal to %d length", len(BaseCurrency), fakeLen))
	}

	if v.Errors.Get("currency") != fmt.Sprintf("The %s must be %d characters long", "currency", fakeLen) {
		t.Error("error: validator error messaging not working")
	}

	v = validator.New(data)
	v.Length("currency", 3)

	if !v.Valid() {
		t.Error("should not have an error, but got one")
	}
}

func TestValidator_Date(t *testing.T) {
	data := make(map[string]string)
	data["from"] = "2022-04-15"
	data["to"] = "2022-04-17"

	v := validator.New(data)
	v.Date("from", "to")

	if !v.Valid() {
		t.Error("got invalid result when dates are valid")
	}

	data = make(map[string]string)
	data["date"] = "20-20-10"
	v = validator.New(data)
	v.Date("date")

	if v.Valid() {
		t.Error("got valid result when date is invalid")
	}
}

func TestValidator_DateInFuture(t *testing.T) {
	data := make(map[string]string)
	data["date"] = "2033-04-15"

	v := validator.New(data)
	v.DateInFuture("date")

	if v.Valid() {
		t.Error("got valid result when date is in the future")
	}

	data = make(map[string]string)
	data["date"] = "2020-04-01"
	v = validator.New(data)
	v.DateInFuture("date")

	if !v.Valid() {
		t.Error("got invalid result when date is not in the future")
	}
}

func TestValidator_ValidRate(t *testing.T) {
	data := make(map[string]string)
	data["rate"] = "1.22"

	v := validator.New(data)
	v.ValidRate("rate")

	if !v.Valid() {
		t.Error("got invalid result when rate is valid")
	}

	data = make(map[string]string)
	data["rate"] = "2020-04-01"
	v = validator.New(data)
	v.ValidRate("rate")

	if v.Valid() {
		t.Error("got valid result when rate is invalid")
	}

	data = make(map[string]string)
	data["rate"] = "-1"
	v = validator.New(data)
	v.ValidRate("rate")

	if v.Valid() {
		t.Error("got valid result when rate is invalid")
	}
}

func TestValidator_NotEqual(t *testing.T) {
	data := make(map[string]string)
	data["currency"] = "usd"

	v := validator.New(data)
	v.NotEqual("currency", BaseCurrency)

	if v.Valid() {
		t.Error("got valid result when field value and comparison value are equal")
	}

	data = make(map[string]string)
	data["currency"] = "chf"
	v = validator.New(data)
	v.NotEqual("currency", BaseCurrency)

	if !v.Valid() {
		t.Error("got invalid result field value and comparison value are different")
	}
}
