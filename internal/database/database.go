package database

import (
	"fmt"

	"github.com/raaj2493/KrishiSetu/internal/buyer"
	"github.com/raaj2493/KrishiSetu/internal/demand"
	"github.com/raaj2493/KrishiSetu/internal/farmer"
	"github.com/raaj2493/KrishiSetu/internal/listing"
	"github.com/raaj2493/KrishiSetu/internal/market"
	"github.com/raaj2493/KrishiSetu/internal/matching"
	"github.com/raaj2493/KrishiSetu/internal/offer"
	"github.com/raaj2493/KrishiSetu/internal/order"

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

	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	return db, nil
}

// migrate ensures the required tables exist so the API works even when the
// external SQL migrations have not all been applied to the target database.
// Before this was added, the `matches` table (and others) could be missing in
// production, causing matching endpoints to fail with internal server errors.
func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&farmer.Farmer{},
		&buyer.Buyer{},
		&listing.CropListing{},
		&demand.Demand{},
		&offer.Offer{},
		&order.Order{},
		&matching.Match{},
		&market.MarketPrice{},
	)
}
