CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,

    offer_id BIGINT NOT NULL UNIQUE
        REFERENCES offers(id),

    listing_id BIGINT NOT NULL
        REFERENCES crop_listings(id),

    buyer_id BIGINT NOT NULL
        REFERENCES buyers(id),

    farmer_id BIGINT NOT NULL
        REFERENCES farmers(id),

    quantity NUMERIC(12, 2) NOT NULL
        CHECK (quantity > 0),

    agreed_price NUMERIC(12, 2) NOT NULL
        CHECK (agreed_price >= 0),

    total_amount NUMERIC(14, 2) NOT NULL
        CHECK (total_amount >= 0),

    status VARCHAR(20) NOT NULL DEFAULT 'CONFIRMED'
        CHECK (
            status IN (
                'CONFIRMED',
                'COMPLETED',
                'CANCELLED'
            )
        ),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_orders_buyer_id
    ON orders(buyer_id);

CREATE INDEX idx_orders_farmer_id
    ON orders(farmer_id);

CREATE INDEX idx_orders_listing_id
    ON orders(listing_id);

CREATE INDEX idx_orders_status
    ON orders(status);