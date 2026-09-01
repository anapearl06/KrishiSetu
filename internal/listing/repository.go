package listing

import "gorm.io/gorm"

type Repository interface {
	Create(l *CropListing) error
	FindByID(id uint) (*CropListing, error)
	FindByFarmerID(farmerID uint) ([]CropListing, error)
	FindAll(crop, state, status string, limit int) ([]CropListing, error)
	Update(l *CropListing) error
	Delete(id uint) error
}

type postgresRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(l *CropListing) error {
	return r.db.Create(l).Error
}

func (r *postgresRepository) FindByID(id uint) (*CropListing, error) {
	var listing CropListing
	if err := r.db.First(&listing, id).Error; err != nil {
		return nil, err
	}
	return &listing, nil
}

func (r *postgresRepository) FindByFarmerID(farmerID uint) ([]CropListing, error) {
	var listings []CropListing
	if err := r.db.Where("farmer_id = ?", farmerID).Order("created_at DESC").Find(&listings).Error; err != nil {
		return nil, err
	}
	return listings, nil
}

func (r *postgresRepository) FindAll(crop, state, status string, limit int) ([]CropListing, error) {
	var listings []CropListing
	query := r.db.Model(&CropListing{})

	if crop != "" {
		query = query.Where("LOWER(crop_name) LIKE ?", "%"+crop+"%")
	}
	if state != "" {
		query = query.Where("LOWER(state) LIKE ?", "%"+state+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query = query.Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&listings).Error; err != nil {
		return nil, err
	}
	return listings, nil
}

func (r *postgresRepository) Update(l *CropListing) error {
	return r.db.Save(l).Error
}

func (r *postgresRepository) Delete(id uint) error {
	return r.db.Delete(&CropListing{}, id).Error
}
