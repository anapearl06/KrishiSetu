CREATE TABLE crop_listings (
    id BIGSERIAL PRIMARY KEY,

    farmer_id BIGINT NOT NULL,

    crop_name VARCHAR(100) NOT NULL,

    quantity DOUBLE PRECISION NOT NULL,

    unit VARCHAR(20) NOT NULL,

    expected_price DOUBLE PRECISION NOT NULL,

    quality_grade VARCHAR(20),

    state VARCHAR(100) NOT NULL,

    district VARCHAR(100) NOT NULL,

    harvest_date TIMESTAMPTZ NOT NULL,

    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_crop_listings_farmer
        FOREIGN KEY (farmer_id)
        REFERENCES farmers(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_crop_listings_farmer_id
    ON crop_listings(farmer_id);

CREATE INDEX idx_crop_listings_crop_name
    ON crop_listings(crop_name);

CREATE INDEX idx_crop_listings_state
    ON crop_listings(state);

CREATE INDEX idx_crop_listings_district
    ON crop_listings(district);

CREATE INDEX idx_crop_listings_status
    ON crop_listings(status);