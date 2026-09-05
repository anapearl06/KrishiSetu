package database

import (
	"fmt"

	"github.com/raaj2493/KrishiSetu/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg config.Config) (*gorm.DB, error) {
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}

	if err := ensureSchema(db); err != nil {
		return nil, fmt.Errorf("ensure schema: %w", err)
	}

	return db, nil
}

// ensureSchema creates tables GORM depends on that are not covered by the
// external SQL migrations on every environment.
//
// Rendering the full schema reports uses GORM AutoMigrate against tables that
// were created by handwritten migrations, which fails at boot because GORM tries
// to reconcile index/constraint names that differ between the two definitions
// (e.g. "uni_farmers_phone"). AutoMigrate is therefore avoided entirely. Only
// tables that may be missing from an existing database are created here, using
// SQL that mirrors the matching migration file so names stay consistent.
//
// The `matches` table (migration 000008) was added after the production
// database had already been provisioned, so it was never created there. Every
// matching endpoint (create, generate, get) then failed with an internal server
// error. This ensures the table exists on startup.
func ensureSchema(db *gorm.DB) error {
	const createMatches = `
CREATE TABLE IF NOT EXISTS matches (
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
`

	const matchesIndexes = `
CREATE INDEX IF NOT EXISTS idx_matches_listing_id
    ON matches(listing_id);

CREATE INDEX IF NOT EXISTS idx_matches_demand_id
    ON matches(demand_id);

CREATE INDEX IF NOT EXISTS idx_matches_score
    ON matches(score DESC);
`

	return db.Exec(createMatches + matchesIndexes).Error
}