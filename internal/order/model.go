package order

import "time"

type Order struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	OfferID      uint      `gorm:"not null;uniqueIndex" json:"offer_id"`
	ListingID    uint      `gorm:"not null;index" json:"listing_id"`
	BuyerID      uint      `gorm:"not null;index" json:"buyer_id"`
	FarmerID     uint      `gorm:"not null;index" json:"farmer_id"`
	Quantity     float64   `gorm:"not null" json:"quantity"`
	AgreedPrice  float64   `gorm:"not null" json:"agreed_price"`
	TotalAmount  float64   `gorm:"not null" json:"total_amount"`
	Status       string    `gorm:"not null;default:CONFIRMED;index" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}