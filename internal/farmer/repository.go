package farmer

import "gorm.io/gorm"

// 1. Define the interface so other packages (like service) can depend on it
type Repository interface {
	Create(f *Farmer) error
}

// 2. Rename the struct to avoid confusion with the interface name
type postgresRepository struct {
	db *gorm.DB
}

// 3. Update constructor to return the interface type
func NewRepository(db *gorm.DB) Repository {
	return &postgresRepository{
		db: db,
	}
}

// 4. Attach the method to the concrete struct
func (r *postgresRepository) Create(f *Farmer) error {
	return r.db.Create(f).Error
}