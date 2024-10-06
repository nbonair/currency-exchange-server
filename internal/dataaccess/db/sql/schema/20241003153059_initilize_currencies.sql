-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS currencies (
    id SERIAL PRIMARY KEY,
    code VARCHAR(3) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    decimal_places INT NOT NULL DEFAULT 2,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Pivot Currencies Table
CREATE TABLE IF NOT EXISTS pivot_currencies (
    id SERIAL PRIMARY KEY,
    currency_id INT REFERENCES currencies(id),
    priority INT NOT NULL DEFAULT 0,
    UNIQUE(currency_id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Exchange Rates Table
CREATE TABLE IF NOT EXISTS exchange_rates (
    base_currency_id INT REFERENCES currencies(id) NOT NULL,
    currency_id INT REFERENCES currencies(id) NOT NULL,
    exchange_rate DECIMAL(20,8) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (base_currency_id, currency_id)
);

-- Exchange Rates History Table
CREATE TABLE IF NOT EXISTS exchange_rate_history (
    id SERIAL PRIMARY KEY,
    base_currency_id INT REFERENCES currencies(id) NOT NULL,
    currency_id INT REFERENCES currencies(id) NOT NULL,
    exchange_rate DECIMAL(20,8) NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ DEFAULT NOW()
);



-- Users Table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE, 
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Rate Subscriptions Table
CREATE TABLE IF NOT EXISTS rate_subscriptions (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    base_currency_id INT REFERENCES currencies(id),
    target_currency_id INT REFERENCES currencies(id),
    pivot_currency_id INT REFERENCES pivot_currencies(id),  
    threshold NUMERIC(10, 4),                
    notification_sent BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes for faster querying
CREATE INDEX idx_from_currency_id ON exchange_rates (source_id);
CREATE INDEX idx_to_currency_id ON exchange_rates (target_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop indexes
DROP INDEX IF EXISTS idx_from_currency_id;
DROP INDEX IF EXISTS idx_to_currency_id;

-- Drop rate_subscriptions table
DROP TABLE IF EXISTS rate_subscriptions;

-- Drop users table
DROP TABLE IF EXISTS users;

-- Drop exchange_rates table
DROP TABLE IF EXISTS exchange_rates;

-- Drop pivot_currencies table
DROP TABLE IF EXISTS pivot_currencies;

-- Drop currencies table
DROP TABLE IF EXISTS currencies;
-- +goose StatementEnd
