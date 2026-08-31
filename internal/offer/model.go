package offer

import "time"

type Offer struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ListingID    uint      `gorm:"not null;index" json:"listing_id"`
	BuyerID      uint      `gorm:"not null;index" json:"buyer_id"`
	Quantity     float64   `gorm:"not null" json:"quantity"`
	OfferedPrice float64   `gorm:"not null" json:"offered_price"`
	Message      string    `json:"message"`
	Status       string    `gorm:"not null;default:PENDING;index" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}