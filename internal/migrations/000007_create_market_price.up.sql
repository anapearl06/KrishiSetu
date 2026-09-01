CREATE TABLE market_prices (
    id BIGSERIAL PRIMARY KEY,

    state VARCHAR(100) NOT NULL,
    district VARCHAR(100) NOT NULL,
    market VARCHAR(150) NOT NULL,

    commodity VARCHAR(150) NOT NULL,
    variety VARCHAR(150),
    grade VARCHAR(100),

    arrival_date DATE NOT NULL,

    min_price DOUBLE PRECISION NOT NULL,
    max_price DOUBLE PRECISION NOT NULL,
    modal_price DOUBLE PRECISION NOT NULL,

    source VARCHAR(100) NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_market_prices_state
    ON market_prices(state);

CREATE INDEX idx_market_prices_district
    ON market_prices(district);

CREATE INDEX idx_market_prices_market
    ON market_prices(market);

CREATE INDEX idx_market_prices_commodity
    ON market_prices(commodity);

CREATE INDEX idx_market_prices_arrival_date
    ON market_prices(arrival_date);