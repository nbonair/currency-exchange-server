-- name: InsertExchangeRateHistory :exec

INSERT INTO exchange_rate_history(base_currency_id, currency_id, exchange_rate, timestamp)
VALUES ($1, $2, $3, $4);
