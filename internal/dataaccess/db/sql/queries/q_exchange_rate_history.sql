-- name: InsertExchangeRateHistory :exec

INSERT INTO exchange_rate_history(base_currency_id, currency_id, exchange_rate, timestamp, created_at)
VALUES (
    (SELECT id FROM currencies c WHERE @base_currency_code = c.code),
    (SELECT id FROM currencies c WHERE @target_currency_code = c.code),
    @exchange_rate,
    NOW(),
    NOW()
);
