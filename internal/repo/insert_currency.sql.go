// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: insert_currency.sql

package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const insertCurrency = `-- name: InsertCurrency :exec
INSERT INTO currencies (code, name, decimal_places, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
ON CONFLICT(code) DO UPDATE 
SET name = EXCLUDED.name, decimal_places = EXCLUDED.decimal_places, updated_at = NOW()
`

type InsertCurrencyParams struct {
	Code          string
	Name          string
	DecimalPlaces int32
}

func (q *Queries) InsertCurrency(ctx context.Context, arg InsertCurrencyParams) error {
	_, err := q.db.Exec(ctx, insertCurrency, arg.Code, arg.Name, arg.DecimalPlaces)
	return err
}

const insertOrUpdatePivotCurrency = `-- name: InsertOrUpdatePivotCurrency :exec
INSERT INTO pivot_currencies (currency_id, priority, created_at, updated_at)
VALUES ($1, $2, NOW(), NOW())
ON CONFLICT (currency_id) DO UPDATE
SET priority = EXCLUDED.priority, updated_at = NOW()
`

type InsertOrUpdatePivotCurrencyParams struct {
	CurrencyID pgtype.Int4
	Priority   int32
}

func (q *Queries) InsertOrUpdatePivotCurrency(ctx context.Context, arg InsertOrUpdatePivotCurrencyParams) error {
	_, err := q.db.Exec(ctx, insertOrUpdatePivotCurrency, arg.CurrencyID, arg.Priority)
	return err
}