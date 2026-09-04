package matching

import (
	"context"
	"errors"
	"testing"

	"github.com/raaj2493/KrishiSetu/internal/demand"
	"github.com/raaj2493/KrishiSetu/internal/listing"
)

// newOwnershipService builds a matching service wired to in-memory
// fakes seeded with:
//   - listing A owned by farmer 1 (Wheat, Aktive)
//   - listing B owned by farmer 3 (Wheat, Aktive)
//   - demand A owned by buyer 2 (Wheat, Aktive)
//   - demand B owned by buyer 4 (Wheat, Aktive)
func newOwnershipService() Service {
	litRepo := &fakeListingRepo{
		listings: map[uint]listing.CropListing{
			1: {
				ID:       1,
				FarmerID: 1,
				CropName: "Wheat",
				Quantity: 100,
				Unit:     "quintal",
				State:    "UP",
				District: "Lucknow",
				Status:   "ACTIVE",
			},
			3: {
				ID:       3,
				FarmerID: 3,
				CropName: "Wheat",
				Quantity: 100,
				Unit:     "quintal",
				State:    "MH",
				District: "Pune",
				Status:   "ACTIVE",
			},
		},
	}

	demRepo := &fakeDemandRepo{
		demands: map[uint]demand.Demand{
			2: {
				ID:       2,
				BuyerID:  2,
				CropName: "Wheat",
				Quantity: 50,
				Unit:     "quintal",
				State:    "UP",
				District: "Lucknow",
				Status:   "ACTIVE",
			},
			4: {
				ID:       4,
				BuyerID:  4,
				CropName: "Wheat",
				Quantity: 50,
				Unit:     "quintal",
				State:    "MH",
				District: "Pune",
				Status:   "ACTIVE",
			},
		},
	}

	matchRepo := &fakeMatchRepo{}

	return NewService(matchRepo, litRepo, demRepo)
}

func TestGenerateMatchesForListingOwnership(t *testing.T) {
	ctx := context.Background()

	t.Run("owner can generate matches", func(t *testing.T) {
		svc := newOwnershipService()
		matches, err := svc.GenerateMatchesForListing(ctx, 1, 1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(matches) == 0 {
			t.Fatal("expected at least one match to be generated")
		}
	})

	t.Run("non-owner cannot generate matches", func(t *testing.T) {
		svc := newOwnershipService()
		_, err := svc.GenerateMatchesForListing(ctx, 3, 1)
		if !errors.Is(err, ErrUnauthorized) {
			t.Fatalf("expected ErrUnauthorized, got %v", err)
		}
	})

	t.Run("missing listing returns not found", func(t *testing.T) {
		svc := newOwnershipService()
		_, err := svc.GenerateMatchesForListing(ctx, 1, 999)
		if !errors.Is(err, ErrListingNotFound) {
			t.Fatalf("expected ErrListingNotFound, got %v", err)
		}
	})
}

func TestGetMatchesForListingOwnership(t *testing.T) {
	ctx := context.Background()

	t.Run("owner can view matches", func(t *testing.T) {
		svc := newOwnershipService()
		if _, err := svc.GenerateMatchesForListing(ctx, 1, 1); err != nil {
			t.Fatalf("setup generate failed: %v", err)
		}
		matches, err := svc.GetMatchesForListing(ctx, 1, 1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(matches) == 0 {
			t.Fatal("expected matches to be visible")
		}
	})

	t.Run("non-owner cannot view matches", func(t *testing.T) {
		svc := newOwnershipService()
		if _, err := svc.GenerateMatchesForListing(ctx, 1, 1); err != nil {
			t.Fatalf("setup generate failed: %v", err)
		}
		_, err := svc.GetMatchesForListing(ctx, 3, 1)
		if !errors.Is(err, ErrUnauthorized) {
			t.Fatalf("expected ErrUnauthorized, got %v", err)
		}
	})

	t.Run("missing listing returns not found", func(t *testing.T) {
		svc := newOwnershipService()
		_, err := svc.GetMatchesForListing(ctx, 1, 999)
		if !errors.Is(err, ErrListingNotFound) {
			t.Fatalf("expected ErrListingNotFound, got %v", err)
		}
	})
}

func TestGenerateMatchesForDemandOwnership(t *testing.T) {
	ctx := context.Background()

	t.Run("owner can generate matches", func(t *testing.T) {
		svc := newOwnershipService()
		matches, err := svc.GenerateMatchesForDemand(ctx, 2, 2)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(matches) == 0 {
			t.Fatal("expected at least one match to be generated")
		}
	})

	t.Run("non-owner cannot generate matches", func(t *testing.T) {
		svc := newOwnershipService()
		_, err := svc.GenerateMatchesForDemand(ctx, 4, 2)
		if !errors.Is(err, ErrUnauthorized) {
			t.Fatalf("expected ErrUnauthorized, got %v", err)
		}
	})

	t.Run("missing demand returns not found", func(t *testing.T) {
		svc := newOwnershipService()
		_, err := svc.GenerateMatchesForDemand(ctx, 2, 999)
		if !errors.Is(err, ErrDemandNotFound) {
			t.Fatalf("expected ErrDemandNotFound, got %v", err)
		}
	})
}

func TestGetMatchesForDemandOwnership(t *testing.T) {
	ctx := context.Background()

	t.Run("owner can view matches", func(t *testing.T) {
		svc := newOwnershipService()
		if _, err := svc.GenerateMatchesForDemand(ctx, 2, 2); err != nil {
			t.Fatalf("setup generate failed: %v", err)
		}
		matches, err := svc.GetMatchesForDemand(ctx, 2, 2)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(matches) == 0 {
			t.Fatal("expected matches to be visible")
		}
	})

	t.Run("non-owner cannot view matches", func(t *testing.T) {
		svc := newOwnershipService()
		if _, err := svc.GenerateMatchesForDemand(ctx, 2, 2); err != nil {
			t.Fatalf("setup generate failed: %v", err)
		}
		_, err := svc.GetMatchesForDemand(ctx, 4, 2)
		if !errors.Is(err, ErrUnauthorized) {
			t.Fatalf("expected ErrUnauthorized, got %v", err)
		}
	})

	t.Run("missing demand returns not found", func(t *testing.T) {
		svc := newOwnershipService()
		_, err := svc.GetMatchesForDemand(ctx, 2, 999)
		if !errors.Is(err, ErrDemandNotFound) {
			t.Fatalf("expected ErrDemandNotFound, got %v", err)
		}
	})
}

func TestCreateMatchOwnership(t *testing.T) {
	ctx := context.Background()

	t.Run("farmer owning listing can create match", func(t *testing.T) {
		svc := newOwnershipService()
		_, err := svc.CreateMatch(ctx, 1, "farmer", 1, 2)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("farmer not owning listing cannot create match", func(t *testing.T) {
		svc := newOwnershipService()
		_, err := svc.CreateMatch(ctx, 3, "farmer", 1, 2)
		if !errors.Is(err, ErrUnauthorized) {
			t.Fatalf("expected ErrUnauthorized, got %v", err)
		}
	})

	t.Run("buyer owning demand can create match", func(t *testing.T) {
		svc := newOwnershipService()
		_, err := svc.CreateMatch(ctx, 2, "buyer", 1, 2)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("buyer not owning demand cannot create match", func(t *testing.T) {
		svc := newOwnershipService()
		_, err := svc.CreateMatch(ctx, 4, "buyer", 1, 2)
		if !errors.Is(err, ErrUnauthorized) {
			t.Fatalf("expected ErrUnauthorized, got %v", err)
		}
	})

	t.Run("farmer cannot create match on another farmers listing", func(t *testing.T) {
		svc := newOwnershipService()
		_, err := svc.CreateMatch(ctx, 1, "farmer", 3, 2)
		if !errors.Is(err, ErrUnauthorized) {
			t.Fatalf("expected ErrUnauthorized, got %v", err)
		}
	})

	t.Run("buyer cannot create match on another buyers demand", func(t *testing.T) {
		svc := newOwnershipService()
		_, err := svc.CreateMatch(ctx, 2, "buyer", 1, 4)
		if !errors.Is(err, ErrUnauthorized) {
			t.Fatalf("expected ErrUnauthorized, got %v", err)
		}
	})

	t.Run("unknown role cannot create match", func(t *testing.T) {
		svc := newOwnershipService()
		_, err := svc.CreateMatch(ctx, 1, "admin", 1, 2)
		if !errors.Is(err, ErrUnauthorized) {
			t.Fatalf("expected ErrUnauthorized, got %v", err)
		}
	})

	t.Run("missing listing returns not found", func(t *testing.T) {
		svc := newOwnershipService()
		_, err := svc.CreateMatch(ctx, 1, "farmer", 999, 2)
		if !errors.Is(err, ErrListingNotFound) {
			t.Fatalf("expected ErrListingNotFound, got %v", err)
		}
	})

	t.Run("missing demand returns not found", func(t *testing.T) {
		svc := newOwnershipService()
		_, err := svc.CreateMatch(ctx, 1, "farmer", 1, 999)
		if !errors.Is(err, ErrDemandNotFound) {
			t.Fatalf("expected ErrDemandNotFound, got %v", err)
		}
	})
}
