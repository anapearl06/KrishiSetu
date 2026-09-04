CREATE TABLE matches (
    id BIGSERIAL PRIMARY KEY,

    listing_id BIGINT NOT NULL,
    demand_id BIGINT NOT NULL,

    score DOUBLE PRECISION NOT NULL,
    level VARCHAR(50) NOT NULL,

    commodity_score DOUBLE PRECISION NOT NULL DEFAULT 0,
    quantity_score DOUBLE PRECISION NOT NULL DEFAULT 0,
    location_score DOUBLE PRECISION NOT NULL DEFAULT 0,
    price_score DOUBLE PRECISION NOT NULL DEFAULT 0,
    grade_score DOUBLE PRECISION NOT NULL DEFAULT 0,

    reasons JSONB NOT NULL DEFAULT '[]'::jsonb,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_matches_listing
        FOREIGN KEY (listing_id)
        REFERENCES crop_listings(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_matches_demand
        FOREIGN KEY (demand_id)
        REFERENCES demands(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_matches_listing_demand
        UNIQUE (listing_id, demand_id)
);

CREATE INDEX idx_matches_listing_id
    ON matches(listing_id);

CREATE INDEX idx_matches_demand_id
    ON matches(demand_id);

CREATE INDEX idx_matches_score
    ON matches(score DESC);