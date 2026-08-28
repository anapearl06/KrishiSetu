package buyer

import "time"

type Buyer struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"not null" json:"name"`
	Phone        string    `gorm:"uniqueIndex;not null" json:"phone"`
	PasswordHash string    `gorm:"not null" json:"-"`
	BusinessName string    `gorm:"not null" json:"business_name"`
	BusinessType string    `json:"business_type"`
	State        string    `gorm:"not null" json:"state"`
	District     string    `gorm:"not null" json:"district"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}