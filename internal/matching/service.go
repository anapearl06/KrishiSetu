package matching

import (
	"context"
	"errors"
	"sort"

	"github.com/raaj2493/KrishiSetu/internal/demand"
	"github.com/raaj2493/KrishiSetu/internal/listing"
	"gorm.io/gorm"
)

// minMatchScore is the minimum score a candidate must reach before it is
// persisted during automatic match generation.
const minMatchScore = 60.0

var (
	ErrListingNotFound = errors.New("listing not found")
	ErrDemandNotFound  = errors.New("demand not found")
	ErrUnauthorized    = errors.New("you do not have access to this resource")
)

type Service interface {
	CreateMatch(
		ctx context.Context,
		userID uint,
		role string,
		listingID uint,
		demandID uint,
	) (*Match, error)

	GenerateMatchesForListing(
		ctx context.Context,
		userID uint,
		listingID uint,
	) ([]Match, error)

	GenerateMatchesForDemand(
		ctx context.Context,
		userID uint,
		demandID uint,
	) ([]Match, error)

	GetMatchesForListing(
		ctx context.Context,
		userID uint,
		listingID uint,
	) ([]Match, error)

	GetMatchesForDemand(
		ctx context.Context,
		userID uint,
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

func (s *service) authenticateListing(
	listingID uint,
	userID uint,
) (*listing.CropListing, error) {
	cropListing, err := s.listingRepo.FindByID(listingID)
	if err != nil {
		return nil, ErrListingNotFound
	}

	if cropListing.FarmerID != userID {
		return nil, ErrUnauthorized
	}

	return cropListing, nil
}

func (s *service) authenticateDemand(
	ctx context.Context,
	demandID uint,
	userID uint,
) (*demand.Demand, error) {
	buyerDemand, err := s.demandRepo.FindByID(ctx, demandID)
	if err != nil {
		return nil, ErrDemandNotFound
	}

	if buyerDemand.BuyerID != userID {
		return nil, ErrUnauthorized
	}

	return buyerDemand, nil
}

func (s *service) CreateMatch(
	ctx context.Context,
	userID uint,
	role string,
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

	// A farmer creating a match must own the listing.
	if role == "farmer" {
		if cropListing.FarmerID != userID {
			return nil, ErrUnauthorized
		}
	}

	// A buyer creating a match must own the demand.
	if role == "buyer" {
		if buyerDemand.BuyerID != userID {
			return nil, ErrUnauthorized
		}
	}

	// Unknown/unexpected role must not be able to create matches.
	if role != "farmer" && role != "buyer" {
		return nil, ErrUnauthorized
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
		// The (listing_id, demand_id) pair already exists. Reuse it
		// rather than failing, keeping generation idempotent.
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			existing, getErr := s.matchRepo.GetByPair(ctx, listingID, demandID)
			if getErr != nil {
				return nil, err
			}
			return existing, nil
		}
		return nil, err
	}

	return match, nil
}

func (s *service) GenerateMatchesForListing(
	ctx context.Context,
	userID uint,
	listingID uint,
) ([]Match, error) {
	cropListing, err := s.authenticateListing(listingID, userID)
	if err != nil {
		return nil, err
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

		if result.Score < minMatchScore {
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
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				existing, getErr := s.matchRepo.GetByPair(
					ctx,
					listingID,
					buyerDemand.ID,
				)
				if getErr != nil {
					return nil, err
				}
				matches = append(matches, *existing)
				continue
			}
			return nil, err
		}

		matches = append(matches, *match)
	}

	sort.SliceStable(matches, func(i, j int) bool {
		return matches[i].Score > matches[j].Score
	})

	return matches, nil
}

func (s *service) GenerateMatchesForDemand(
	ctx context.Context,
	userID uint,
	demandID uint,
) ([]Match, error) {
	buyerDemand, err := s.authenticateDemand(ctx, demandID, userID)
	if err != nil {
		return nil, err
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

		if result.Score < minMatchScore {
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
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				existing, getErr := s.matchRepo.GetByPair(
					ctx,
					cropListing.ID,
					demandID,
				)
				if getErr != nil {
					return nil, err
				}
				matches = append(matches, *existing)
				continue
			}
			return nil, err
		}

		matches = append(matches, *match)
	}

	sort.SliceStable(matches, func(i, j int) bool {
		return matches[i].Score > matches[j].Score
	})

	return matches, nil
}

func (s *service) GetMatchesForListing(
	ctx context.Context,
	userID uint,
	listingID uint,
) ([]Match, error) {
	if _, err := s.authenticateListing(listingID, userID); err != nil {
		return nil, err
	}

	return s.matchRepo.GetByListing(
		ctx,
		listingID,
	)
}

func (s *service) GetMatchesForDemand(
	ctx context.Context,
	userID uint,
	demandID uint,
) ([]Match, error) {
	if _, err := s.authenticateDemand(ctx, demandID, userID); err != nil {
		return nil, err
	}

	return s.matchRepo.GetByDemand(
		ctx,
		demandID,
	)
}
