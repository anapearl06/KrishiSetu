package demand

import (
	"context"
	"strings"
	"testing"
	"time"
)

// listFakeRepo mirrors the production List filtering (crop, state, district,
// status, price range) so the untested ListDemands service path can be
// exercised deterministically.
type listFakeRepo struct {
	demands []Demand
}

func (f *listFakeRepo) Create(ctx context.Context, d *Demand) error {
	f.demands = append(f.demands, *d)
	return nil
}

func (f *listFakeRepo) FindByID(ctx context.Context, id uint) (*Demand, error) {
	for _, d := range f.demands {
		if d.ID == id {
			return &d, nil
		}
	}
	return nil, ErrDemandNotFound
}

func (f *listFakeRepo) FindByBuyer(ctx context.Context, buyerID uint) ([]Demand, error) {
	var out []Demand
	for _, d := range f.demands {
		if d.BuyerID == buyerID {
			out = append(out, d)
		}
	}
	return out, nil
}

func (f *listFakeRepo) Update(ctx context.Context, d *Demand) error {
	return nil
}

func (f *listFakeRepo) Cancel(ctx context.Context, id uint) error {
	return nil
}

func (f *listFakeRepo) List(ctx context.Context, filters DemandFilters) ([]Demand, error) {
	var out []Demand
	for _, d := range f.demands {
		if filters.CropName != "" && !strings.EqualFold(d.CropName, filters.CropName) {
			continue
		}
		if filters.State != "" && !strings.EqualFold(d.State, filters.State) {
			continue
		}
		if filters.District != "" && !strings.EqualFold(d.District, filters.District) {
			continue
		}
		if filters.Status != "" && d.Status != filters.Status {
			continue
		}
		if filters.MinPrice != nil && d.TargetPrice < *filters.MinPrice {
			continue
		}
		if filters.MaxPrice != nil && d.TargetPrice > *filters.MaxPrice {
			continue
		}
		out = append(out, d)
	}
	return out, nil
}

func newListFakeRepo() *listFakeRepo {
	return &listFakeRepo{}
}

func seedDemands(ctx context.Context, repo *listFakeRepo) {
	produce := []struct {
		crop, state, district, status string
		price                         float64
	}{
		{"Wheat", "Uttar Pradesh", "Agra", "ACTIVE", 2000},
		{"wheat", "Uttar Pradesh", "Lucknow", "ACTIVE", 1900},
		{"Rice", "Punjab", "Ludhiana", "ACTIVE", 2200},
		{"Onion", "Maharashtra", "Nashik", "COMPLETED", 800},
	}
	for i, p := range produce {
		repo.Create(ctx, &Demand{
			ID:          uint(i + 1),
			BuyerID:     10,
			CropName:    p.crop,
			State:       p.state,
			District:    p.district,
			Status:      p.status,
			TargetPrice: p.price,
			RequiredBy:  time.Now().AddDate(0, 0, 7),
		})
	}
}

func TestListDemands_DefaultsToActive(t *testing.T) {
	ctx := context.Background()
	repo := newListFakeRepo()
	svc := NewService(repo)
	seedDemands(ctx, repo)

	demands, err := svc.ListDemands(ctx, DemandFilters{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Only ACTIVE demands (3 of 4) should be returned because Status defaults
	// to ACTIVE when left empty.
	if len(demands) != 3 {
		t.Fatalf("expected 3 active demands, got %d", len(demands))
	}
	for _, d := range demands {
		if d.Status != "ACTIVE" {
			t.Errorf("expected only ACTIVE demands, got %s", d.Status)
		}
	}
}

func TestListDemands_RespectsExplicitStatusFilter(t *testing.T) {
	ctx := context.Background()
	repo := newListFakeRepo()
	svc := NewService(repo)
	seedDemands(ctx, repo)

	demands, err := svc.ListDemands(ctx, DemandFilters{Status: "COMPLETED"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(demands) != 1 || demands[0].CropName != "Onion" {
		t.Fatalf("expected only the COMPLETED Onion demand, got %+v", demands)
	}
}

func TestListDemands_FiltersByCropCaseInsensitive(t *testing.T) {
	ctx := context.Background()
	repo := newListFakeRepo()
	svc := NewService(repo)
	seedDemands(ctx, repo)

	demands, err := svc.ListDemands(ctx, DemandFilters{CropName: "WHEAT", Status: "ACTIVE"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// "WHEAT" matches both "Wheat" and "wheat".
	if len(demands) != 2 {
		t.Fatalf("expected 2 wheat demands, got %d", len(demands))
	}
}

func TestListDemands_FiltersByDistrict(t *testing.T) {
	ctx := context.Background()
	repo := newListFakeRepo()
	svc := NewService(repo)
	seedDemands(ctx, repo)

	demands, err := svc.ListDemands(ctx, DemandFilters{District: "Agra", Status: "ACTIVE"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(demands) != 1 || demands[0].District != "Agra" {
		t.Fatalf("expected only Agra demand, got %+v", demands)
	}
}

func TestListDemands_FiltersByPriceRange(t *testing.T) {
	ctx := context.Background()
	repo := newListFakeRepo()
	svc := NewService(repo)
	seedDemands(ctx, repo)

	min, max := 1900.0, 2000.0
	demands, err := svc.ListDemands(ctx, DemandFilters{MinPrice: &min, MaxPrice: &max, Status: "ACTIVE"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// Wheat 2000 and wheat 1900 fall in range; Rice 2200 excluded.
	if len(demands) != 2 {
		t.Fatalf("expected 2 demands in price range, got %d", len(demands))
	}
}

func TestListDemands_NoMatchReturnsEmpty(t *testing.T) {
	ctx := context.Background()
	repo := newListFakeRepo()
	svc := NewService(repo)
	seedDemands(ctx, repo)

	demands, err := svc.ListDemands(ctx, DemandFilters{CropName: "Tomato", Status: "ACTIVE"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(demands) != 0 {
		t.Fatalf("expected empty result, got %d", len(demands))
	}
}

func TestCreateDemand_InvalidPrice(t *testing.T) {
	repo := newListFakeRepo()
	svc := NewService(repo)
	ctx := context.Background()

	d := &Demand{CropName: "Wheat", Quantity: 100, Unit: "kg", TargetPrice: 0, State: "UP", District: "Agra", RequiredBy: time.Now().AddDate(0, 0, 7)}
	if err := svc.CreateDemand(ctx, 10, d); err != ErrInvalidPrice {
		t.Errorf("expected ErrInvalidPrice, got %v", err)
	}
}

func TestCreateDemand_InvalidUnit(t *testing.T) {
	repo := newListFakeRepo()
	svc := NewService(repo)
	ctx := context.Background()

	d := &Demand{CropName: "Wheat", Quantity: 100, Unit: "", TargetPrice: 2000, State: "UP", District: "Agra", RequiredBy: time.Now().AddDate(0, 0, 7)}
	if err := svc.CreateDemand(ctx, 10, d); err != ErrInvalidUnit {
		t.Errorf("expected ErrInvalidUnit, got %v", err)
	}
}

func TestCreateDemand_InvalidLocation(t *testing.T) {
	repo := newListFakeRepo()
	svc := NewService(repo)
	ctx := context.Background()

	d := &Demand{CropName: "Wheat", Quantity: 100, Unit: "kg", TargetPrice: 2000, State: "", District: "", RequiredBy: time.Now().AddDate(0, 0, 7)}
	if err := svc.CreateDemand(ctx, 10, d); err != ErrInvalidLocation {
		t.Errorf("expected ErrInvalidLocation, got %v", err)
	}
}

func TestCreateDemand_InvalidRequiredBy(t *testing.T) {
	repo := newListFakeRepo()
	svc := NewService(repo)
	ctx := context.Background()

	d := &Demand{CropName: "Wheat", Quantity: 100, Unit: "kg", TargetPrice: 2000, State: "UP", District: "Agra"}
	if err := svc.CreateDemand(ctx, 10, d); err != ErrInvalidRequiredBy {
		t.Errorf("expected ErrInvalidRequiredBy, got %v", err)
	}
}

func TestUpdateDemand_InvalidPrice(t *testing.T) {
	repo := newListFakeRepo()
	svc := NewService(repo)
	ctx := context.Background()

	repo.Create(ctx, &Demand{ID: 1, BuyerID: 10, CropName: "Wheat", Status: "ACTIVE", RequiredBy: time.Now().AddDate(0, 0, 7)})

	update := &Demand{ID: 1, CropName: "Wheat", Quantity: 100, Unit: "kg", TargetPrice: 0, State: "UP", District: "Agra", RequiredBy: time.Now().AddDate(0, 0, 7)}
	if err := svc.UpdateDemand(ctx, 10, update); err != ErrInvalidPrice {
		t.Errorf("expected ErrInvalidPrice, got %v", err)
	}
}
