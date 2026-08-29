package demand

import (
	"context"
	"errors"
	"strings"
)

var (
	ErrDemandNotFound  = errors.New("demand not found")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrDemandNotActive = errors.New("demand is not active")

	ErrInvalidQuantity = errors.New(
		"quantity must be greater than zero",
	)

	ErrInvalidPrice = errors.New(
		"target price must be greater than zero",
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

	ErrInvalidRequiredBy = errors.New(
		"required_by date is required",
	)
)

type Service interface {
	CreateDemand(
		ctx context.Context,
		buyerID uint,
		demand *Demand,
	) error

	GetDemand(
		ctx context.Context,
		id uint,
	) (*Demand, error)

	GetBuyerDemands(
		ctx context.Context,
		buyerID uint,
	) ([]Demand, error)

	UpdateDemand(
		ctx context.Context,
		buyerID uint,
		demand *Demand,
	) error

	CancelDemand(
		ctx context.Context,
		buyerID uint,
		demandID uint,
	) error

	ListDemands(
		ctx context.Context,
		filters DemandFilters,
	) ([]Demand, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateDemand(
	ctx context.Context,
	buyerID uint,
	demand *Demand,
) error {
	if buyerID == 0 {
		return ErrUnauthorized
	}

	if err := validateDemand(demand); err != nil {
		return err
	}

	demand.BuyerID = buyerID
	demand.Status = "ACTIVE"

	return s.repo.Create(ctx, demand)
}

func (s *service) GetDemand(
	ctx context.Context,
	id uint,
) (*Demand, error) {
	demand, err := s.repo.FindByID(ctx, id)

	if err != nil {
		return nil, ErrDemandNotFound
	}

	return demand, nil
}

func (s *service) GetBuyerDemands(
	ctx context.Context,
	buyerID uint,
) ([]Demand, error) {
	if buyerID == 0 {
		return nil, ErrUnauthorized
	}

	return s.repo.FindByBuyer(ctx, buyerID)
}

func (s *service) UpdateDemand(
	ctx context.Context,
	buyerID uint,
	demand *Demand,
) error {
	existing, err := s.repo.FindByID(ctx, demand.ID)

	if err != nil {
		return ErrDemandNotFound
	}

	if existing.BuyerID != buyerID {
		return ErrUnauthorized
	}

	if existing.Status != "ACTIVE" {
		return ErrDemandNotActive
	}

	if err := validateDemand(demand); err != nil {
		return err
	}

	// Never trust ownership or status from the client.
	demand.BuyerID = existing.BuyerID
	demand.Status = existing.Status

	return s.repo.Update(ctx, demand)
}

func (s *service) CancelDemand(
	ctx context.Context,
	buyerID uint,
	demandID uint,
) error {
	demand, err := s.repo.FindByID(ctx, demandID)

	if err != nil {
		return ErrDemandNotFound
	}

	if demand.BuyerID != buyerID {
		return ErrUnauthorized
	}

	if demand.Status != "ACTIVE" {
		return ErrDemandNotActive
	}

	return s.repo.Cancel(ctx, demandID)
}

func (s *service) ListDemands(
	ctx context.Context,
	filters DemandFilters,
) ([]Demand, error) {
	if filters.Status == "" {
		filters.Status = "ACTIVE"
	}

	return s.repo.List(ctx, filters)
}

func validateDemand(demand *Demand) error {
	if strings.TrimSpace(demand.CropName) == "" {
		return ErrInvalidCropName
	}

	if demand.Quantity <= 0 {
		return ErrInvalidQuantity
	}

	if demand.TargetPrice <= 0 {
		return ErrInvalidPrice
	}

	if strings.TrimSpace(demand.Unit) == "" {
		return ErrInvalidUnit
	}

	if strings.TrimSpace(demand.State) == "" ||
		strings.TrimSpace(demand.District) == "" {
		return ErrInvalidLocation
	}

	if demand.RequiredBy.IsZero() {
		return ErrInvalidRequiredBy
	}

	return nil
}