package farmer

import "gorm.io/gorm"

// 1. Define the interface
type Repository interface {
	Create(f *Farmer) error
	FindByPhone(phone string) (*Farmer, error)
}

// 2. Define the concrete struct
type postgresRepository struct {
	db *gorm.DB
}

// 3. Constructor returns the Repository interface
func NewRepository(db *gorm.DB) Repository {
	return &postgresRepository{
		db: db,
	}
}

// 4. Implement Create on *postgresRepository
func (r *postgresRepository) Create(f *Farmer) error {
	return r.db.Create(f).Error
}

// 5. Implement FindByPhone on *postgresRepository (fix the receiver name here)
func (r *postgresRepository) FindByPhone(phone string) (*Farmer, error) {
	var farmer Farmer

	// Note: Make sure your column name matches your migration (e.g., phone_number vs phone)
	if err := r.db.Where("phone = ?", phone).First(&farmer).Error; err != nil {
		return nil, err
	}

	return &farmer, nil
}
