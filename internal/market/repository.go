package market

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, price *MarketPrice) error

	GetLatestByCommodity(
		ctx context.Context,
		commodity string,
	) (*MarketPrice, error)

	GetByCommodityAndLocation(
		ctx context.Context,
		commodity string,
		state string,
		district string,
	) ([]MarketPrice, error)

	GetHistoricalByCommodity(
		ctx context.Context,
		commodity string,
		from time.Time,
		to time.Time,
	) ([]MarketPrice, error)
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
	price *MarketPrice,
) error {
	return r.db.WithContext(ctx).
		Create(price).
		Error
}

func (r *repository) GetLatestByCommodity(
	ctx context.Context,
	commodity string,
) (*MarketPrice, error) {
	var price MarketPrice

	err := r.db.WithContext(ctx).
		Where("commodity = ?", commodity).
		Order("arrival_date DESC").
		First(&price).
		Error

	if err != nil {
		return nil, err
	}

	return &price, nil
}

func (r *repository) GetByCommodityAndLocation(
	ctx context.Context,
	commodity string,
	state string,
	district string,
) ([]MarketPrice, error) {
	var prices []MarketPrice

	query := r.db.WithContext(ctx).
		Where("commodity = ?", commodity)

	if state != "" {
		query = query.Where("state = ?", state)
	}

	if district != "" {
		query = query.Where("district = ?", district)
	}

	err := query.
		Order("arrival_date DESC").
		Find(&prices).
		Error

	return prices, err
}

func (r *repository) GetHistoricalByCommodity(
	ctx context.Context,
	commodity string,
	from time.Time,
	to time.Time,
) ([]MarketPrice, error) {
	var prices []MarketPrice

	err := r.db.WithContext(ctx).
		Where("commodity = ?", commodity).
		Where("arrival_date >= ?", from).
		Where("arrival_date <= ?", to).
		Order("arrival_date ASC").
		Find(&prices).
		Error

	return prices, err
}