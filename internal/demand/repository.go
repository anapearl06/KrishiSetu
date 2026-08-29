package demand

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, demand *Demand) error
	FindByID(ctx context.Context, id uint) (*Demand, error)
	FindByBuyer(ctx context.Context, buyerID uint) ([]Demand, error)
	Update(ctx context.Context, demand *Demand) error
	Cancel(ctx context.Context, id uint) error
	List(ctx context.Context, filters DemandFilters) ([]Demand, error)
}

type DemandFilters struct {
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

func (r *repository) Create(
	ctx context.Context,
	demand *Demand,
) error {
	return r.db.WithContext(ctx).
		Create(demand).
		Error
}

func (r *repository) FindByID(
	ctx context.Context,
	id uint,
) (*Demand, error) {
	var demand Demand

	err := r.db.WithContext(ctx).
		First(&demand, id).
		Error

	if err != nil {
		return nil, err
	}

	return &demand, nil
}

func (r *repository) FindByBuyer(
	ctx context.Context,
	buyerID uint,
) ([]Demand, error) {
	var demands []Demand

	err := r.db.WithContext(ctx).
		Where("buyer_id = ?", buyerID).
		Order("created_at DESC").
		Find(&demands).
		Error

	return demands, err
}

func (r *repository) Update(
	ctx context.Context,
	demand *Demand,
) error {
	return r.db.WithContext(ctx).
		Save(demand).
		Error
}

func (r *repository) Cancel(
	ctx context.Context,
	id uint,
) error {
	return r.db.WithContext(ctx).
		Model(&Demand{}).
		Where("id = ?", id).
		Update("status", "CANCELLED").
		Error
}

func (r *repository) List(
	ctx context.Context,
	filters DemandFilters,
) ([]Demand, error) {
	var demands []Demand

	query := r.db.WithContext(ctx).
		Model(&Demand{})

	if filters.CropName != "" {
		query = query.Where(
			"LOWER(crop_name) = LOWER(?)",
			filters.CropName,
		)
	}

	if filters.State != "" {
		query = query.Where(
			"LOWER(state) = LOWER(?)",
			filters.State,
		)
	}

	if filters.District != "" {
		query = query.Where(
			"LOWER(district) = LOWER(?)",
			filters.District,
		)
	}

	if filters.Status != "" {
		query = query.Where(
			"status = ?",
			filters.Status,
		)
	}

	if filters.MinPrice != nil {
		query = query.Where(
			"target_price >= ?",
			*filters.MinPrice,
		)
	}

	if filters.MaxPrice != nil {
		query = query.Where(
			"target_price <= ?",
			*filters.MaxPrice,
		)
	}

	err := query.
		Order("created_at DESC").
		Find(&demands).
		Error

	return demands, err
}