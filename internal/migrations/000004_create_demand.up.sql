CREATE TABLE demands (
    id BIGSERIAL PRIMARY KEY,

    buyer_id BIGINT NOT NULL,

    crop_name VARCHAR(100) NOT NULL,

    quantity DOUBLE PRECISION NOT NULL,

    unit VARCHAR(20) NOT NULL,

    target_price DOUBLE PRECISION NOT NULL,

    state VARCHAR(100) NOT NULL,

    district VARCHAR(100) NOT NULL,

    required_by TIMESTAMPTZ NOT NULL,

    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_demands_buyer
        FOREIGN KEY (buyer_id)
        REFERENCES buyers(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_demands_buyer_id
    ON demands(buyer_id);

CREATE INDEX idx_demands_crop_name
    ON demands(crop_name);

CREATE INDEX idx_demands_state
    ON demands(state);

CREATE INDEX idx_demands_district
    ON demands(district);

CREATE INDEX idx_demands_status
    ON demands(status);


    