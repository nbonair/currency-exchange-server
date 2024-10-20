-- name: InsertExchangeRateHistory :exec

INSERT INTO exchange_rate_history(base_currency_id, currency_id, exchange_rate, timestamp, created_at)
VALUES (
    (SELECT id FROM currencies c WHERE @base_currency_code = c.code),
    (SELECT id FROM currencies c WHERE @target_currency_code = c.code),
    @exchange_rate,
    NOW(),
    NOW()
);

-- name: GetExchangeRateHistory :many

SELECT
    erh.exchange_rate AS rate,
    erh.timestamp,
    bc.code as base_currency_code,
    tc.code as target_currency_code
FROM exchange_rate_history erh
JOIN currencies bc ON erh.base_currency_id = bc.id
JOIN currencies tc ON erh.base_currency_id = tc.id
WHERE bc.code = @base_currency_code
    AND tc.code = @target_currency_code
    AND erh.timestamp BETWEEN @start_time AND @end_time
ORDER BY erh.timestamp ASC;


