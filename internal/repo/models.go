// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package repo

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Currency struct {
	ID            int32
	Code          string
	Name          string
	DecimalPlaces int32
	CreatedAt     pgtype.Timestamptz
	UpdatedAt     pgtype.Timestamptz
}

type ExchangeRate struct {
	ID           int32
	SourceID     int32
	TargetID     int32
	PivotID      int32
	ExchangeRate pgtype.Numeric
	Priority     int32
	Timestamp    pgtype.Timestamptz
	CreatedAt    pgtype.Timestamptz
	UpdatedAt    pgtype.Timestamptz
}

type PivotCurrency struct {
	ID         int32
	CurrencyID pgtype.Int4
	Priority   int32
	CreatedAt  pgtype.Timestamptz
	UpdatedAt  pgtype.Timestamptz
}

type RateSubscription struct {
	ID               int32
	UserID           pgtype.Int4
	BaseCurrencyID   pgtype.Int4
	TargetCurrencyID pgtype.Int4
	PivotCurrencyID  pgtype.Int4
	Threshold        pgtype.Numeric
	NotificationSent pgtype.Bool
	CreatedAt        pgtype.Timestamptz
	UpdatedAt        pgtype.Timestamptz
}

type User struct {
	ID        int32
	Email     string
	CreatedAt pgtype.Timestamptz
}