package matching

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, match *Match) error

	GetByListing(
		ctx context.Context,
		listingID uint,
	) ([]Match, error)

	GetByDemand(
		ctx context.Context,
		demandID uint,
	) ([]Match, error)

	GetBestMatchesForListing(
		ctx context.Context,
		listingID uint,
	) ([]Match, error)

	GetBestMatchesForDemand(
		ctx context.Context,
		demandID uint,
	) ([]Match, error)

	GetByPair(
		ctx context.Context,
		listingID uint,
		demandID uint,
	) (*Match, error)
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
	match *Match,
) error {
	return r.db.
		WithContext(ctx).
		Create(match).
		Error
}

func (r *repository) GetByListing(
	ctx context.Context,
	listingID uint,
) ([]Match, error) {
	var matches []Match

	err := r.db.
		WithContext(ctx).
		Where("listing_id = ?", listingID).
		Order("score DESC").
		Find(&matches).
		Error

	return matches, err
}

func (r *repository) GetByDemand(
	ctx context.Context,
	demandID uint,
) ([]Match, error) {
	var matches []Match

	err := r.db.
		WithContext(ctx).
		Where("demand_id = ?", demandID).
		Order("score DESC").
		Find(&matches).
		Error

	return matches, err
}

func (r *repository) GetBestMatchesForListing(
	ctx context.Context,
	listingID uint,
) ([]Match, error) {
	var matches []Match

	err := r.db.
		WithContext(ctx).
		Where("listing_id = ?", listingID).
		Where("score >= ?", 60).
		Order("score DESC").
		Find(&matches).
		Error

	return matches, err
}

func (r *repository) GetBestMatchesForDemand(
	ctx context.Context,
	demandID uint,
) ([]Match, error) {
	var matches []Match

	err := r.db.
		WithContext(ctx).
		Where("demand_id = ?", demandID).
		Where("score >= ?", 60).
		Order("score DESC").
		Find(&matches).
		Error

	return matches, err
}

func (r *repository) GetByPair(
	ctx context.Context,
	listingID uint,
	demandID uint,
) (*Match, error) {
	var match Match

	err := r.db.
		WithContext(ctx).
		Where("listing_id = ?", listingID).
		Where("demand_id = ?", demandID).
		First(&match).
		Error

	if err != nil {
		return nil, err
	}

	return &match, nil
}