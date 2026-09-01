CREATE TABLE offers (
    id BIGSERIAL PRIMARY KEY,

    listing_id BIGINT NOT NULL,
    farmer_id BIGINT NOT NULL,
    buyer_id BIGINT NOT NULL,

    price DOUBLE PRECISION NOT NULL,
    quantity DOUBLE PRECISION NOT NULL,
    message TEXT,

    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_offers_listing
        FOREIGN KEY (listing_id)
        REFERENCES crop_listings(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_offers_farmer
        FOREIGN KEY (farmer_id)
        REFERENCES farmers(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_offers_buyer
        FOREIGN KEY (buyer_id)
        REFERENCES buyers(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_offers_listing_id ON offers(listing_id);
CREATE INDEX idx_offers_farmer_id ON offers(farmer_id);
CREATE INDEX idx_offers_buyer_id ON offers(buyer_id);
CREATE INDEX idx_offers_status ON offers(status);
