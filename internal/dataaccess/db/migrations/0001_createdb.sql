-- +goose Up
CREATE TABLE currencies (
    id SERIAL PRIMARY KEY,
    code VARCHAR(3) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE rates (
    id SERIAL PRIMARY KEY,
    base_currency_id INT REFERENCES currencies(id),
    target_currency_id INT REFERENCES currencies(id),
    rate NUMERIC(10, 4) NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS rates;
DROP TABLE IF EXISTS currencies;
