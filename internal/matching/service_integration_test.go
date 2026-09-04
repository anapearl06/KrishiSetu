package matching

import (
	"context"
	"errors"
	"testing"

	"github.com/raaj2493/KrishiSetu/internal/demand"
	"github.com/raaj2493/KrishiSetu/internal/listing"
)

// intgFixture bundles the three in-memory repositories so tests can
// inspect persisted records and seed data.
type intgFixture struct {
	svc     Service
	listing *fakeListingRepo
	demand  *fakeDemandRepo
	match   *fakeMatchRepo
}

// newIntegrationService returns a service wired to fresh in-memory fakes.
func newIntegrationService() *intgFixture {
	lit := &fakeListingRepo{listings: map[uint]listing.CropListing{}}
	dem := &fakeDemandRepo{demands: map[uint]demand.Demand{}}
	mat := &fakeMatchRepo{}

	return &intgFixture{
		svc:     NewService(mat, lit, dem),
		listing: lit,
		demand:  dem,
		match:   mat,
	}
}

// seedListing inserts a listing owned by the given farmer.
func (f *intgFixture) seedListing(farmerID, id uint, crop string, qty float64, unit string, price float64) {
	f.listing.listings[id] = listing.CropListing{
		ID:            id,
		FarmerID:      farmerID,
		CropName:      crop,
		Quantity:      qty,
		Unit:          unit,
		ExpectedPrice: price,
		QualityGrade:  "A",
		State:         "UP",
		District:      "Lucknow",
		Status:        "ACTIVE",
	}
}

// seedDemand inserts a demand owned by the given buyer.
func (f *intgFixture) seedDemand(buyerID, id uint, crop string, qty float64, unit string, price float64) {
	f.demand.demands[id] = demand.Demand{
		ID:          id,
		BuyerID:     buyerID,
		CropName:    crop,
		Quantity:    qty,
		Unit:        unit,
		TargetPrice: price,
		State:       "UP",
		District:    "Lucknow",
		Status:      "ACTIVE",
	}
}

func (f *intgFixture) matchCount() int {
	return len(f.match.matches)
}

func TestIntegrationFarmerFlow(t *testing.T) {
	ctx := context.Background()
	f := newIntegrationService()
	f.seedListing(1, 10, "Wheat", 100, "quintal", 2500)
	f.seedDemand(2, 20, "Wheat", 50, "quintal", 2500)

	matches, err := f.svc.GenerateMatchesForListing(ctx, 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	if matches[0].ListingID != 10 || matches[0].DemandID != 20 {
		t.Fatalf("match references wrong pair: listing=%d demand=%d", matches[0].ListingID, matches[0].DemandID)
	}
	if matches[0].Score < minMatchScore {
		t.Fatalf("expected score >= %v, got %v", minMatchScore, matches[0].Score)
	}
	if matches[0].Reasons == nil || len(matches[0].Reasons) == 0 {
		t.Fatal("expected explainable reasons to be present")
	}
}

func TestIntegrationBuyerFlow(t *testing.T) {
	ctx := context.Background()
	f := newIntegrationService()
	f.seedDemand(2, 20, "Wheat", 50, "quintal", 2500)
	f.seedListing(1, 10, "Wheat", 100, "quintal", 2500)

	matches, err := f.svc.GenerateMatchesForDemand(ctx, 2, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	if matches[0].ListingID != 10 || matches[0].DemandID != 20 {
		t.Fatalf("match references wrong pair")
	}
}

func TestIntegrationIncompatibleCommodity(t *testing.T) {
	ctx := context.Background()
	f := newIntegrationService()
	f.seedListing(1, 10, "Wheat", 100, "quintal", 2500)
	f.seedDemand(2, 20, "Rice", 50, "quintal", 2500)

	matches, err := f.svc.GenerateMatchesForListing(ctx, 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(matches) != 0 {
		t.Fatalf("expected no match for incompatible commodity, got %d", len(matches))
	}
	if f.matchCount() != 0 {
		t.Fatalf("expected nothing persisted, got %d records", f.matchCount())
	}
}

func TestIntegrationLowScoreNotPersisted(t *testing.T) {
	ctx := context.Background()
	f := newIntegrationService()
	// Matching commodity but different units, a different state, and a far
	// price keeps the overall score below the 60 threshold.
	f.seedListing(1, 10, "Wheat", 100, "kg", 100)
	f.demand.demands[20] = demand.Demand{
		ID: 20, BuyerID: 2, CropName: "Wheat", Quantity: 5000,
		Unit: "quintal", TargetPrice: 10000, State: "MH", District: "Pune",
		Status: "ACTIVE",
	}

	matches, err := f.svc.GenerateMatchesForListing(ctx, 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(matches) != 0 {
		t.Fatalf("expected no match below threshold, got %d", len(matches))
	}
	if f.matchCount() != 0 {
		t.Fatalf("expected low-score match to be filtered, got %d records", f.matchCount())
	}
}

func TestIntegrationInactiveExcluded(t *testing.T) {
	ctx := context.Background()
	f := newIntegrationService()
	f.seedListing(1, 10, "Wheat", 100, "quintal", 2500)
	f.seedDemand(2, 20, "Wheat", 50, "quintal", 2500)
	// Mark the demand inactive.
	f.demand.demands[20] = demand.Demand{
		ID: 20, BuyerID: 2, CropName: "Wheat", Quantity: 50,
		Unit: "quintal", TargetPrice: 2500, State: "UP", District: "Lucknow",
		Status: "CLOSED",
	}

	matches, err := f.svc.GenerateMatchesForListing(ctx, 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(matches) != 0 {
		t.Fatalf("expected inactive demand to be excluded, got %d", len(matches))
	}
}

func TestIntegrationNoCandidates(t *testing.T) {
	ctx := context.Background()
	f := newIntegrationService()
	f.seedListing(1, 10, "Wheat", 100, "quintal", 2500)

	matches, err := f.svc.GenerateMatchesForListing(ctx, 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(matches) != 0 {
		t.Fatalf("expected no matches with no demands, got %d", len(matches))
	}
}

func TestIntegrationDuplicateGeneration(t *testing.T) {
	ctx := context.Background()
	f := newIntegrationService()
	f.seedListing(1, 10, "Wheat", 100, "quintal", 2500)
	f.seedDemand(2, 20, "Wheat", 50, "quintal", 2500)

	if _, err := f.svc.GenerateMatchesForListing(ctx, 1, 10); err != nil {
		t.Fatalf("first generation failed: %v", err)
	}
	first := f.matchCount()

	matches, err := f.svc.GenerateMatchesForListing(ctx, 1, 10)
	if err != nil {
		t.Fatalf("second generation failed: %v", err)
	}
	if f.matchCount() != first {
		t.Fatalf("expected no duplicate records: before=%d after=%d", first, f.matchCount())
	}
	if len(matches) != 1 {
		t.Fatalf("expected 1 match returned, got %d", len(matches))
	}
}

func TestIntegrationDuplicateGenerationBuyerSide(t *testing.T) {
	ctx := context.Background()
	f := newIntegrationService()
	f.seedDemand(2, 20, "Wheat", 50, "quintal", 2500)
	f.seedListing(1, 10, "Wheat", 100, "quintal", 2500)

	if _, err := f.svc.GenerateMatchesForDemand(ctx, 2, 20); err != nil {
		t.Fatalf("first generation failed: %v", err)
	}
	first := f.matchCount()

	matches, err := f.svc.GenerateMatchesForDemand(ctx, 2, 20)
	if err != nil {
		t.Fatalf("second generation failed: %v", err)
	}
	if f.matchCount() != first {
		t.Fatalf("expected no duplicate records: before=%d after=%d", first, f.matchCount())
	}
	if len(matches) != 1 {
		t.Fatalf("expected 1 match returned, got %d", len(matches))
	}
}

func TestIntegrationUniqueConstraintRace(t *testing.T) {
	ctx := context.Background()
	f := newIntegrationService()
	f.seedListing(1, 10, "Wheat", 100, "quintal", 2500)
	f.seedDemand(2, 20, "Wheat", 50, "quintal", 2500)
	// Simulate a pre-existing record so a fresh Create collides with the
	// unique constraint.
	f.match.matches = append(f.match.matches, Match{
		ListingID: 10, DemandID: 20, Score: 85, Level: "Good Match",
	})
	f.match.emulateUniqueConstraint = true

	matches, err := f.svc.GenerateMatchesForListing(ctx, 1, 10)
	if err != nil {
		t.Fatalf("expected duplicate collision to be handled, got error: %v", err)
	}
	if len(matches) != 1 {
		t.Fatalf("expected the existing match to be returned, got %d", len(matches))
	}
	if matches[0].ListingID != 10 || matches[0].DemandID != 20 {
		t.Fatalf("returned wrong existing match")
	}
	if f.matchCount() != 1 {
		t.Fatalf("expected no duplicate record appended, got %d", f.matchCount())
	}
}

func TestIntegrationExistingPairReuseInCreateMatch(t *testing.T) {
	ctx := context.Background()
	f := newIntegrationService()
	f.seedListing(1, 10, "Wheat", 100, "quintal", 2500)
	f.seedDemand(2, 20, "Wheat", 50, "quintal", 2500)
	f.match.matches = append(f.match.matches, Match{
		ListingID: 10, DemandID: 20, Score: 85, Level: "Good Match",
	})
	f.match.emulateUniqueConstraint = true

	created, err := f.svc.CreateMatch(ctx, 1, "farmer", 10, 20)
	if err != nil {
		t.Fatalf("expected collision handled in CreateMatch, got error: %v", err)
	}
	if created.ListingID != 10 || created.DemandID != 20 {
		t.Fatalf("returned wrong existing match")
	}
	if f.matchCount() != 1 {
		t.Fatalf("expected no duplicate insert, got %d", f.matchCount())
	}
}

func TestIntegrationSortedByScore(t *testing.T) {
	ctx := context.Background()
	f := newIntegrationService()
	// Same commodity but differing prices produce distinguishable scores.
	f.seedListing(1, 10, "Wheat", 100, "quintal", 2500)
	f.seedDemand(2, 21, "Wheat", 50, "quintal", 2500)  // exact price -> high score
	f.seedDemand(2, 22, "Wheat", 50, "quintal", 2800)  // ~12% -> lower score
	f.seedDemand(2, 20, "Wheat", 50, "quintal", 10000) // far -> weak but commodity matches

	matches, err := f.svc.GenerateMatchesForListing(ctx, 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(matches) < 3 {
		t.Fatalf("expected at least 3 matches, got %d", len(matches))
	}
	for i := 1; i < len(matches); i++ {
		if matches[i-1].Score < matches[i].Score {
			t.Fatalf("matches not sorted by score DESC at index %d: %v then %v", i, matches[i-1].Score, matches[i].Score)
		}
	}
}

func TestIntegrationMissingResourceIDs(t *testing.T) {
	ctx := context.Background()
	f := newIntegrationService()

	if _, err := f.svc.GenerateMatchesForListing(ctx, 1, 999); !errors.Is(err, ErrListingNotFound) {
		t.Fatalf("expected ErrListingNotFound for missing listing, got %v", err)
	}
	if _, err := f.svc.GenerateMatchesForDemand(ctx, 2, 999); !errors.Is(err, ErrDemandNotFound) {
		t.Fatalf("expected ErrDemandNotFound for missing demand, got %v", err)
	}
}
