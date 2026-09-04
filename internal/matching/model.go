package matching

import "time"

type Match struct {
	ID uint `gorm:"primaryKey" json:"id"`

	ListingID uint `gorm:"not null;index;uniqueIndex:idx_listing_demand" json:"listing_id"`
	DemandID  uint `gorm:"not null;index;uniqueIndex:idx_listing_demand" json:"demand_id"`

	Score float64 `gorm:"not null;index" json:"score"`
	Level string  `gorm:"not null" json:"level"`

	CommodityScore float64 `gorm:"not null;default:0" json:"commodity_score"`
	QuantityScore  float64 `gorm:"not null;default:0" json:"quantity_score"`
	LocationScore  float64 `gorm:"not null;default:0" json:"location_score"`
	PriceScore     float64 `gorm:"not null;default:0" json:"price_score"`
	GradeScore     float64 `gorm:"not null;default:0" json:"grade_score"`

	Reasons []string `gorm:"type:jsonb;serializer:json;not null" json:"reasons"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}