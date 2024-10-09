// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"time"
)

type Currency struct {
	ID            int32
	Code          string
	Name          string
	DecimalPlaces int32
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type ExchangeRate struct {
	BaseCurrencyID int32
	CurrencyID     int32
	ExchangeRate   float64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type ExchangeRateHistory struct {
	ID             int32
	BaseCurrencyID int32
	CurrencyID     int32
	ExchangeRate   float64
	Timestamp      time.Time
	CreatedAt      time.Time
}

type PivotCurrency struct {
	ID         int32
	CurrencyID int32
	Priority   int32
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
