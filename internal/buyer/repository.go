package buyer

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, buyer *Buyer) error
	FindByID(ctx context.Context, id uint) (*Buyer, error)
	FindByPhone(ctx context.Context, phone string) (*Buyer, error)
	Update(ctx context.Context, buyer *Buyer) error
}

type postgresRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &postgresRepository{
		db: db,
	}
}

func (r *postgresRepository) Create(ctx context.Context, buyer *Buyer) error {
	return r.db.WithContext(ctx).Create(buyer).Error
}

func (r *postgresRepository) FindByID(ctx context.Context, id uint) (*Buyer, error) {
	var buyer Buyer

	if err := r.db.WithContext(ctx).First(&buyer, id).Error; err != nil {
		return nil, err
	}

	return &buyer, nil
}

func (r *postgresRepository) FindByPhone(ctx context.Context, phone string) (*Buyer, error) {
	var buyer Buyer

	if err := r.db.WithContext(ctx).
		Where("phone = ?", phone).
		First(&buyer).Error; err != nil {
		return nil, err
	}

	return &buyer, nil
}

func (r *postgresRepository) Update(ctx context.Context, buyer *Buyer) error {
	return r.db.WithContext(ctx).Save(buyer).Error
}