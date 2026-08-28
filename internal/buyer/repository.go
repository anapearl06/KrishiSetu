package buyer

import (
	"context"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, buyer *Buyer) error {
	return r.db.WithContext(ctx).Create(buyer).Error
}

func (r *Repository) FindByID(ctx context.Context, id uint) (*Buyer, error) {
	var buyer Buyer

	err := r.db.WithContext(ctx).
		First(&buyer, id).
		Error

	if err != nil {
		return nil, err
	}

	return &buyer, nil
}

func (r *Repository) FindByPhone(ctx context.Context, phone string) (*Buyer, error) {
	var buyer Buyer

	err := r.db.WithContext(ctx).
		Where("phone = ?", phone).
		First(&buyer).
		Error

	if err != nil {
		return nil, err
	}

	return &buyer, nil
}

func (r *Repository) Update(ctx context.Context, buyer *Buyer) error {
	return r.db.WithContext(ctx).Save(buyer).Error
}