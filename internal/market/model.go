package market

import "time"

type MarketPrice struct {
	ID uint `gorm:"primaryKey"`

	State     string `gorm:"not null;index"`
	District  string `gorm:"not null;index"`
	Market    string `gorm:"not null;index"`
	Commodity string `gorm:"not null;index"`
	Variety   string
	Grade     string

	ArrivalDate time.Time `gorm:"not null;index"`

	MinPrice   float64 `gorm:"not null"`
	MaxPrice   float64 `gorm:"not null"`
	ModalPrice float64 `gorm:"not null"`

	Source    string    `gorm:"not null"`
	CreatedAt time.Time
}