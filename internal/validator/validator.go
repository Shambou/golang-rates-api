package validator

import (
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// Validator - creates a custom data struct, embeds mux.Vars
type Validator struct {
	Data   map[string]string
	Errors errors
}

// Valid - returns true if there are no errors, otherwise false
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// New - initializes a form struct
func New(data map[string]string) *Validator {
	return &Validator{
		data,
		errors(map[string][]string{}),
	}
}

// Length - checks length of field
func (v *Validator) Length(field string, length int) {
	value := v.Get(field)
	if len(value) != length {
		v.Errors.Add(field, fmt.Sprintf("The %s must be %d characters long", field, length))
	}
}

// Date - checks if field is valid date
func (v *Validator) Date(fields ...string) {
	for _, field := range fields {
		value := v.Get(field)
		_, err := time.Parse("2006-01-02", value)
		if err != nil {
			v.Errors.Add(field, fmt.Sprintf("%s date is invalid", value))
		}
	}
}

// ValidRate - checks if field value is valid decimal first then checks if it's less than or equal to zero
func (v *Validator) ValidRate(field string) {
	value, err := decimal.NewFromString(v.Get(field))
	if err != nil {
		v.Errors.Add(field, "The %s is invalid")
	}

	if value.LessThanOrEqual(decimal.Zero) {
		v.Errors.Add(field, fmt.Sprintf("The %s is invalid", field))
	}
}

// NotEqualTo - checks if field value is not equal to comparison value
func (v *Validator) NotEqualTo(field string, comparisonValue string) {
	value := v.Get(field)

	if strings.ToTitle(value) == strings.ToTitle(comparisonValue) {
		v.Errors.Add(field, fmt.Sprintf("The %s can't be %s", field, comparisonValue))
	}
}

func (v *Validator) Get(key string) string {
	vs := v.Data[key]
	if len(vs) == 0 {
		return ""
	}

	return vs
}
