-- +goose Up
-- +goose StatementBegin

-- Currencies Table
CREATE TABLE IF NOT EXISTS currencies (
    id SERIAL PRIMARY KEY,
    code VARCHAR(3) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    decimal_places INT NOT NULL DEFAULT 2,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Pivot Currencies Table
CREATE TABLE IF NOT EXISTS pivot_currencies (
    id SERIAL PRIMARY KEY,
    currency_id INT NOT NULL REFERENCES currencies(id) ON DELETE CASCADE,
    priority INT NOT NULL DEFAULT 0,
    UNIQUE(currency_id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Exchange Rates Table
CREATE TABLE IF NOT EXISTS exchange_rates (
    base_currency_id INT NOT NULL REFERENCES currencies(id) ON DELETE CASCADE,
    currency_id INT NOT NULL REFERENCES currencies(id) ON DELETE CASCADE,
    exchange_rate FLOAT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (base_currency_id, currency_id)
);

-- Exchange Rates History Table
CREATE TABLE IF NOT EXISTS exchange_rate_history (
    id SERIAL PRIMARY KEY,
    base_currency_id INT NOT NULL REFERENCES currencies(id) ON DELETE CASCADE,
    currency_id INT NOT NULL REFERENCES currencies(id) ON DELETE CASCADE,
    exchange_rate FLOAT NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);



-- Indexes for faster querying
CREATE INDEX idx_exchange_rates_base_currency_id ON exchange_rates (base_currency_id);
CREATE INDEX idx_exchange_rates_currency_id ON exchange_rates (currency_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop indexes
DROP INDEX IF EXISTS idx_exchange_rates_base_currency_id;
DROP INDEX IF EXISTS idx_exchange_rates_currency_id;

-- Drop tables
DROP TABLE IF EXISTS exchange_rate_history;
DROP TABLE IF EXISTS exchange_rates;
DROP TABLE IF EXISTS pivot_currencies;
DROP TABLE IF EXISTS currencies;

-- +goose StatementEnd
