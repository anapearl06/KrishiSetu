package listing

import (
	"context"
	"errors"
	"strings"
)

var (
	ErrListingNotFound  = errors.New("listing not found")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrListingNotActive = errors.New("listing is not active")

	ErrInvalidQuantity = errors.New(
		"quantity must be greater than zero",
	)

	ErrInvalidPrice = errors.New(
		"expected price must be greater than zero",
	)

	ErrInvalidCropName = errors.New(
		"crop name is required",
	)

	ErrInvalidUnit = errors.New(
		"unit is required",
	)

	ErrInvalidLocation = errors.New(
		"state and district are required",
	)

	ErrInvalidHarvestDate = errors.New(
		"harvest date is required",
	)
)

type Service interface {
	CreateListing(
		ctx context.Context,
		farmerID uint,
		listing *CropListing,
	) error

	GetListing(
		ctx context.Context,
		id uint,
	) (*CropListing, error)

	GetFarmerListings(
		ctx context.Context,
		farmerID uint,
	) ([]CropListing, error)

	UpdateListing(
		ctx context.Context,
		farmerID uint,
		listing *CropListing,
	) error

	CancelListing(
		ctx context.Context,
		farmerID uint,
		listingID uint,
	) error

	ListListings(
		ctx context.Context,
		filters ListingFilters,
	) ([]CropListing, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateListing(
	ctx context.Context,
	farmerID uint,
	listing *CropListing,
) error {
	if farmerID == 0 {
		return ErrUnauthorized
	}

	if err := validateListing(listing); err != nil {
		return err
	}

	listing.FarmerID = farmerID
	listing.Status = "ACTIVE"

	return s.repo.Create(ctx, listing)
}

func (s *service) GetListing(
	ctx context.Context,
	id uint,
) (*CropListing, error) {
	listing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrListingNotFound
	}

	return listing, nil
}

func (s *service) GetFarmerListings(
	ctx context.Context,
	farmerID uint,
) ([]CropListing, error) {
	if farmerID == 0 {
		return nil, ErrUnauthorized
	}

	return s.repo.FindByFarmer(ctx, farmerID)
}

func (s *service) UpdateListing(
	ctx context.Context,
	farmerID uint,
	listing *CropListing,
) error {
	existing, err := s.repo.FindByID(ctx, listing.ID)
	if err != nil {
		return ErrListingNotFound
	}

	if existing.FarmerID != farmerID {
		return ErrUnauthorized
	}

	if existing.Status != "ACTIVE" {
		return ErrListingNotActive
	}

	if err := validateListing(listing); err != nil {
		return err
	}

	// Never trust ownership or status received from the client.
	listing.FarmerID = existing.FarmerID
	listing.Status = existing.Status

	return s.repo.Update(ctx, listing)
}

func (s *service) CancelListing(
	ctx context.Context,
	farmerID uint,
	listingID uint,
) error {
	listing, err := s.repo.FindByID(ctx, listingID)
	if err != nil {
		return ErrListingNotFound
	}

	if listing.FarmerID != farmerID {
		return ErrUnauthorized
	}

	if listing.Status != "ACTIVE" {
		return ErrListingNotActive
	}

	return s.repo.Cancel(ctx, listingID)
}

func (s *service) ListListings(
	ctx context.Context,
	filters ListingFilters,
) ([]CropListing, error) {
	// Marketplace discovery should show active supply by default.
	if filters.Status == "" {
		filters.Status = "ACTIVE"
	}

	return s.repo.List(ctx, filters)
}

func validateListing(listing *CropListing) error {
	if strings.TrimSpace(listing.CropName) == "" {
		return ErrInvalidCropName
	}

	if listing.Quantity <= 0 {
		return ErrInvalidQuantity
	}

	if listing.ExpectedPrice <= 0 {
		return ErrInvalidPrice
	}

	if strings.TrimSpace(listing.Unit) == "" {
		return ErrInvalidUnit
	}

	if strings.TrimSpace(listing.State) == "" ||
		strings.TrimSpace(listing.District) == "" {
		return ErrInvalidLocation
	}

	if listing.HarvestDate.IsZero() {
		return ErrInvalidHarvestDate
	}

	return nil
}