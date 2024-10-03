-- name: InsertCurrency :exec
INSERT INTO currencies (code, name, decimal_places, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
ON CONFLICT(code) DO UPDATE 
SET name = EXCLUDED.name, decimal_places = EXCLUDED.decimal_places, updated_at = NOW();

-- name: InsertOrUpdatePivotCurrency :exec
INSERT INTO pivot_currencies (currency_id, priority, created_at, updated_at)
VALUES ($1, $2, NOW(), NOW())
ON CONFLICT (currency_id) DO UPDATE
SET priority = EXCLUDED.priority, updated_at = NOW();