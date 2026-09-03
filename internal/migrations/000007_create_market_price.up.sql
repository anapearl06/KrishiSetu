CREATE TABLE market_prices (
    id BIGSERIAL PRIMARY KEY,

    state VARCHAR(100) NOT NULL,
    district VARCHAR(100) NOT NULL,
    market VARCHAR(150) NOT NULL,
    commodity VARCHAR(100) NOT NULL,
    variety VARCHAR(100) NOT NULL,
    grade VARCHAR(100) NOT NULL,

    arrival_date DATE NOT NULL,

    min_price NUMERIC(10, 2) NOT NULL,
    max_price NUMERIC(10, 2) NOT NULL,
    modal_price NUMERIC(10, 2) NOT NULL,

    source VARCHAR(255) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_market_prices_commodity
    ON market_prices(commodity);

CREATE INDEX idx_market_prices_state
    ON market_prices(state);

CREATE INDEX idx_market_prices_district
    ON market_prices(district);

CREATE INDEX idx_market_prices_market
    ON market_prices(market);

CREATE INDEX idx_market_prices_arrival_date
    ON market_prices(arrival_date);

CREATE INDEX idx_market_prices_commodity_state
    ON market_prices(commodity, state);

CREATE INDEX idx_market_prices_commodity_date
    ON market_prices(commodity, arrival_date);