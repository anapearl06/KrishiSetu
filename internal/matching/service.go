package matching

import (
	"context"
	"errors"

	"github.com/raaj2493/KrishiSetu/internal/demand"
	"github.com/raaj2493/KrishiSetu/internal/listing"
)

var (
	ErrListingNotFound = errors.New("listing not found")
	ErrDemandNotFound  = errors.New("demand not found")
)

type Service interface {
	CreateMatch(
		ctx context.Context,
		listingID uint,
		demandID uint,
	) (*Match, error)

	GenerateMatchesForListing(
		ctx context.Context,
		listingID uint,
	) ([]Match, error)

	GenerateMatchesForDemand(
		ctx context.Context,
		demandID uint,
	) ([]Match, error)

	GetMatchesForListing(
		ctx context.Context,
		listingID uint,
	) ([]Match, error)

	GetMatchesForDemand(
		ctx context.Context,
		demandID uint,
	) ([]Match, error)
}

type service struct {
	matchRepo   Repository
	listingRepo listing.Repository
	demandRepo  demand.Repository
}

func NewService(
	matchRepo Repository,
	listingRepo listing.Repository,
	demandRepo demand.Repository,
) Service {
	return &service{
		matchRepo:   matchRepo,
		listingRepo: listingRepo,
		demandRepo:  demandRepo,
	}
}

func (s *service) CreateMatch(
	ctx context.Context,
	listingID uint,
	demandID uint,
) (*Match, error) {
	cropListing, err := s.listingRepo.FindByID(listingID)
	if err != nil {
		return nil, ErrListingNotFound
	}

	buyerDemand, err := s.demandRepo.FindByID(ctx, demandID)
	if err != nil {
		return nil, ErrDemandNotFound
	}

	result := CalculateMatch(
		*cropListing,
		*buyerDemand,
	)

	match := &Match{
		ListingID: listingID,
		DemandID:  demandID,

		Score: result.Score,
		Level: result.Level,

		CommodityScore: result.CommodityScore,
		QuantityScore:  result.QuantityScore,
		LocationScore:  result.LocationScore,
		PriceScore:     result.PriceScore,
		GradeScore:     result.GradeScore,

		Reasons: result.Reasons,
	}

	if err := s.matchRepo.Create(ctx, match); err != nil {
		return nil, err
	}

	return match, nil
}

func (s *service) GenerateMatchesForListing(
	ctx context.Context,
	listingID uint,
) ([]Match, error) {
	cropListing, err := s.listingRepo.FindByID(listingID)
	if err != nil {
		return nil, ErrListingNotFound
	}

	demands, err := s.demandRepo.List(
		ctx,
		demand.DemandFilters{
			CropName: cropListing.CropName,
			Status:   "ACTIVE",
		},
	)
	if err != nil {
		return nil, err
	}

	matches := make([]Match, 0)

	for i := range demands {
		buyerDemand := demands[i]

		result := CalculateMatch(
			*cropListing,
			buyerDemand,
		)

		if result.Score < 60 {
			continue
		}

		existing, err := s.matchRepo.GetByPair(
			ctx,
			listingID,
			buyerDemand.ID,
		)

		if err == nil && existing != nil {
			matches = append(matches, *existing)
			continue
		}

		match := &Match{
			ListingID: listingID,
			DemandID:  buyerDemand.ID,

			Score: result.Score,
			Level: result.Level,

			CommodityScore: result.CommodityScore,
			QuantityScore:  result.QuantityScore,
			LocationScore:  result.LocationScore,
			PriceScore:     result.PriceScore,
			GradeScore:     result.GradeScore,

			Reasons: result.Reasons,
		}

		if err := s.matchRepo.Create(ctx, match); err != nil {
			return nil, err
		}

		matches = append(matches, *match)
	}

	return matches, nil
}

func (s *service) GenerateMatchesForDemand(
	ctx context.Context,
	demandID uint,
) ([]Match, error) {
	buyerDemand, err := s.demandRepo.FindByID(
		ctx,
		demandID,
	)
	if err != nil {
		return nil, ErrDemandNotFound
	}

	listings, err := s.listingRepo.FindAll(
		buyerDemand.CropName,
		"",
		"ACTIVE",
		0,
	)
	if err != nil {
		return nil, err
	}

	matches := make([]Match, 0)

	for i := range listings {
		cropListing := listings[i]

		result := CalculateMatch(
			cropListing,
			*buyerDemand,
		)

		if result.Score < 60 {
			continue
		}

		existing, err := s.matchRepo.GetByPair(
			ctx,
			cropListing.ID,
			demandID,
		)

		if err == nil && existing != nil {
			matches = append(matches, *existing)
			continue
		}

		match := &Match{
			ListingID: cropListing.ID,
			DemandID:  demandID,

			Score: result.Score,
			Level: result.Level,

			CommodityScore: result.CommodityScore,
			QuantityScore:  result.QuantityScore,
			LocationScore:  result.LocationScore,
			PriceScore:     result.PriceScore,
			GradeScore:     result.GradeScore,

			Reasons: result.Reasons,
		}

		if err := s.matchRepo.Create(ctx, match); err != nil {
			return nil, err
		}

		matches = append(matches, *match)
	}

	return matches, nil
}

func (s *service) GetMatchesForListing(
	ctx context.Context,
	listingID uint,
) ([]Match, error) {
	return s.matchRepo.GetByListing(
		ctx,
		listingID,
	)
}

func (s *service) GetMatchesForDemand(
	ctx context.Context,
	demandID uint,
) ([]Match, error) {
	return s.matchRepo.GetByDemand(
		ctx,
		demandID,
	)
}