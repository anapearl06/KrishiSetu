package offer

import "time"

type Offer struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ListingID uint      `gorm:"not null;index" json:"listing_id"`
	FarmerID  uint      `gorm:"not null;index" json:"farmer_id"`
	BuyerID   uint      `gorm:"not null;index" json:"buyer_id"`
	Price     float64   `gorm:"not null" json:"price"`
	Quantity  float64   `gorm:"not null" json:"quantity"`
	Message   string    `json:"message"`
	Status    string    `gorm:"not null;default:'PENDING'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
