-- name: UpsertExchangeRate :exec
INSERT INTO exchange_rates(base_currency_id, currency_id, exchange_rate, updated_at)
VALUES ($1, $2, $3, $4)
ON CONFLICT (base_currency_id, currency_id)
DO UPDATE SET exchange_rate = EXCLUDED.exchange_rate, updated_at = EXCLUDED.updated_at;

-- name: GetExchangeRate :one
SELECT 
    er.exchange_rate,
    er.created_at,
    er.updated_at,
    bc.code as base_currency_code,
    tc.code as target_currency_code
FROM exchange_rates er 
JOIN currencies bc ON er.base_currency_id = bc.id
JOIN currencies tc ON er.currency_id = tc.id
WHERE bc.code = @base_currency_code  AND tc.code = @target_currency_code;

