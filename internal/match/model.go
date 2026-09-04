package match

import "time"

type Match struct {
    ID uint `gorm:"primaryKey"`

    ListingID uint `gorm:"not null;index"`
    DemandID  uint `gorm:"not null;index"`

    Score float64 `gorm:"not null"`
    Level string  `gorm:"not null"`

    CommodityScore float64 `gorm:"not null"`
    QuantityScore  float64 `gorm:"not null"`
    LocationScore  float64 `gorm:"not null"`
    PriceScore     float64 `gorm:"not null"`
    GradeScore     float64 `gorm:"not null"`

    Reasons []string `gorm:"type:jsonb"`

    CreatedAt time.Time
    UpdatedAt time.Time
}