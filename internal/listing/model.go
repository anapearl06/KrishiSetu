package listing

import "time"

type CropListing struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	FarmerID      uint      `gorm:"not null;index" json:"farmer_id"`
	CropName      string    `gorm:"not null;index" json:"crop_name"`
	Quantity      float64   `gorm:"not null" json:"quantity"`
	Unit          string    `gorm:"not null" json:"unit"`
	ExpectedPrice float64   `gorm:"not null" json:"expected_price"`
	QualityGrade  string    `json:"quality_grade"`
	State         string    `gorm:"not null;index" json:"state"`
	District      string    `gorm:"not null;index" json:"district"`
	HarvestDate   time.Time `gorm:"not null" json:"harvest_date"`
	Status        string    `gorm:"not null;index" json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}