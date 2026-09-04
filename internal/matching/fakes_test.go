package matching

import (
	"context"
	"errors"

	"github.com/raaj2493/KrishiSetu/internal/demand"
	"github.com/raaj2493/KrishiSetu/internal/listing"
	"gorm.io/gorm"
)

// fakeListingRepo is an in-memory listing.Repository used for
// service-level ownership tests. No database is required.
type fakeListingRepo struct {
	listings map[uint]listing.CropListing
}

func (f *fakeListingRepo) Create(l *listing.CropListing) error {
	f.listings[l.ID] = *l
	return nil
}

func (f *fakeListingRepo) FindByID(id uint) (*listing.CropListing, error) {
	l, ok := f.listings[id]
	if !ok {
		return nil, errors.New("not found")
	}
	copy := l
	return &copy, nil
}

func (f *fakeListingRepo) FindByFarmerID(farmerID uint) ([]listing.CropListing, error) {
	var result []listing.CropListing
	for _, l := range f.listings {
		if l.FarmerID == farmerID {
			result = append(result, l)
		}
	}
	return result, nil
}

func (f *fakeListingRepo) FindAll(crop, state, status string, limit int) ([]listing.CropListing, error) {
	var result []listing.CropListing
	for _, l := range f.listings {
		if status != "" && l.Status != status {
			continue
		}
		result = append(result, l)
	}
	return result, nil
}

func (f *fakeListingRepo) Update(l *listing.CropListing) error {
	f.listings[l.ID] = *l
	return nil
}

func (f *fakeListingRepo) Delete(id uint) error {
	delete(f.listings, id)
	return nil
}

// fakeDemandRepo is an in-memory demand.Repository for tests.
type fakeDemandRepo struct {
	demands map[uint]demand.Demand
}

func (f *fakeDemandRepo) Create(ctx context.Context, d *demand.Demand) error {
	f.demands[d.ID] = *d
	return nil
}

func (f *fakeDemandRepo) FindByID(ctx context.Context, id uint) (*demand.Demand, error) {
	d, ok := f.demands[id]
	if !ok {
		return nil, errors.New("not found")
	}
	copy := d
	return &copy, nil
}

func (f *fakeDemandRepo) FindByBuyer(ctx context.Context, buyerID uint) ([]demand.Demand, error) {
	var result []demand.Demand
	for _, d := range f.demands {
		if d.BuyerID == buyerID {
			result = append(result, d)
		}
	}
	return result, nil
}

func (f *fakeDemandRepo) Update(ctx context.Context, d *demand.Demand) error {
	f.demands[d.ID] = *d
	return nil
}

func (f *fakeDemandRepo) Cancel(ctx context.Context, id uint) error {
	delete(f.demands, id)
	return nil
}

func (f *fakeDemandRepo) List(ctx context.Context, filters demand.DemandFilters) ([]demand.Demand, error) {
	var result []demand.Demand
	for _, d := range f.demands {
		if filters.Status != "" && d.Status != filters.Status {
			continue
		}
		if filters.CropName != "" && d.CropName != filters.CropName {
			continue
		}
		result = append(result, d)
	}
	return result, nil
}

// fakeMatchRepo is an in-memory matching.Repository for tests.
type fakeMatchRepo struct {
	matches []Match

	// emulateUniqueConstraint, when true, makes Create fail with a
	// duplicate-key error if the (listing_id, demand_id) pair already
	// exists — mirroring the real database unique constraint.
	emulateUniqueConstraint bool
}

func (f *fakeMatchRepo) Create(ctx context.Context, m *Match) error {
	if f.emulateUniqueConstraint {
		for _, existing := range f.matches {
			if existing.ListingID == m.ListingID && existing.DemandID == m.DemandID {
				return gorm.ErrDuplicatedKey
			}
		}
	}
	f.matches = append(f.matches, *m)
	return nil
}

func (f *fakeMatchRepo) GetByListing(ctx context.Context, listingID uint) ([]Match, error) {
	var result []Match
	for _, m := range f.matches {
		if m.ListingID == listingID {
			result = append(result, m)
		}
	}
	return result, nil
}

func (f *fakeMatchRepo) GetByDemand(ctx context.Context, demandID uint) ([]Match, error) {
	var result []Match
	for _, m := range f.matches {
		if m.DemandID == demandID {
			result = append(result, m)
		}
	}
	return result, nil
}

func (f *fakeMatchRepo) GetBestMatchesForListing(ctx context.Context, listingID uint) ([]Match, error) {
	return f.GetByListing(ctx, listingID)
}

func (f *fakeMatchRepo) GetBestMatchesForDemand(ctx context.Context, demandID uint) ([]Match, error) {
	return f.GetByDemand(ctx, demandID)
}

func (f *fakeMatchRepo) GetByPair(ctx context.Context, listingID, demandID uint) (*Match, error) {
	for _, m := range f.matches {
		if m.ListingID == listingID && m.DemandID == demandID {
			copy := m
			return &copy, nil
		}
	}
	return nil, errors.New("not found")
}
