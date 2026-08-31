CREATE TABLE offers (
    id BIGSERIAL PRIMARY KEY,

    listing_id BIGINT NOT NULL REFERENCES crop_listings(id) ON DELETE CASCADE,
    buyer_id BIGINT NOT NULL REFERENCES buyers(id) ON DELETE CASCADE,

    quantity NUMERIC(12, 2) NOT NULL CHECK (quantity > 0),
    offered_price NUMERIC(12, 2) NOT NULL CHECK (offered_price >= 0),

    message TEXT,

    status VARCHAR(20) NOT NULL DEFAULT 'PENDING'
        CHECK (status IN ('PENDING', 'ACCEPTED', 'REJECTED', 'CANCELLED')),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_offers_listing_id ON offers(listing_id);
CREATE INDEX idx_offers_buyer_id ON offers(buyer_id);
CREATE INDEX idx_offers_status ON offers(status);