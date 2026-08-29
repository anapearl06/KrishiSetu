package demand

import "time"

type Demand struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	BuyerID     uint      `gorm:"not null;index" json:"buyer_id"`
	CropName    string    `gorm:"not null;index" json:"crop_name"`
	Quantity    float64   `gorm:"not null" json:"quantity"`
	Unit        string    `gorm:"not null" json:"unit"`
	TargetPrice float64   `gorm:"not null" json:"target_price"`
	State       string    `gorm:"not null;index" json:"state"`
	District    string    `gorm:"not null;index" json:"district"`
	RequiredBy  time.Time `gorm:"not null" json:"required_by"`
	Status      string    `gorm:"not null;index" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}