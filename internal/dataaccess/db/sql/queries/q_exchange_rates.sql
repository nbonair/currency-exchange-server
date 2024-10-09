-- name: UpsertExchangeRate :exec
INSERT INTO exchange_rates(base_currency_id, currency_id, exchange_rate, updated_at)
VALUES ($1, $2, $3, $4)
ON CONFLICT (base_currency_id, currency_id)
DO UPDATE SET exchange_rate = EXCLUDED.exchange_rate, updated_at = EXCLUDED.updated_at;



