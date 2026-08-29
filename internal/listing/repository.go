package listing

import (
	"context"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, listing *CropListing) error
	FindByID(ctx context.Context, id uint) (*CropListing, error)
	FindByFarmer(ctx context.Context, farmerID uint) ([]CropListing, error)
	Update(ctx context.Context, listing *CropListing) error
	Cancel(ctx context.Context, id uint) error
	List(ctx context.Context, filters ListingFilters) ([]CropListing, error)
}

type ListingFilters struct {
	CropName string
	State    string
	District string
	Status   string
	MinPrice *float64
	MaxPrice *float64
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context , listing *CropListing) error {
	return r.db.WithContext(ctx).Create(listing).Error
}


func (r *repository) FindByID(ctx context.Context, id uint) (*CropListing, error) {
	var listing CropListing

	err := r.db.WithContext(ctx).
		First(&listing, id).Error

	if err != nil {
		return nil, err
	}

	return &listing, nil
}



func (r *repository) FindByFarmer(ctx context.Context, farmerID uint) ([]CropListing, error) {
	var listings []CropListing

	err := r.db.WithContext(ctx).
		Where("farmer_id = ?", farmerID).
		Order("created_at DESC").
		Find(&listings).Error

	return listings, err
}


func (r *repository) Update(ctx context.Context, listing *CropListing) error {
	return r.db.WithContext(ctx).
		Save(listing).Error
}


func (r *repository) Cancel(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&CropListing{}).
		Where("id = ?", id).
		Update("status", "CANCELLED").Error
}


func (r *repository) List(ctx context.Context, filters ListingFilters) ([]CropListing, error) {
	var listings []CropListing

	query := r.db.WithContext(ctx).
		Model(&CropListing{})

	if filters.CropName != "" {
		query = query.Where("LOWER(crop_name) = LOWER(?)", filters.CropName)
	}

	if filters.State != "" {
		query = query.Where("LOWER(state) = LOWER(?)", filters.State)
	}

	if filters.District != "" {
		query = query.Where("LOWER(district) = LOWER(?)", filters.District)
	}

	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}

	if filters.MinPrice != nil {
		query = query.Where("expected_price >= ?", *filters.MinPrice)
	}

	if filters.MaxPrice != nil {
		query = query.Where("expected_price <= ?", *filters.MaxPrice)
	}

	err := query.
		Order("created_at DESC").
		Find(&listings).Error

	return listings, err
}





