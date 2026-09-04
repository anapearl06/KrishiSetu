package demand

import (
	"context"
	"testing"
	"time"
)

type fakeDemandRepo struct {
	demands map[uint]Demand
}

func (f *fakeDemandRepo) Create(ctx context.Context, d *Demand) error {
	f.demands[d.ID] = *d
	return nil
}
func (f *fakeDemandRepo) FindByID(ctx context.Context, id uint) (*Demand, error) {
	d, ok := f.demands[id]
	if !ok {
		return nil, ErrDemandNotFound
	}
	return &d, nil
}
func (f *fakeDemandRepo) FindByBuyer(ctx context.Context, buyerID uint) ([]Demand, error) {
	var result []Demand
	for _, d := range f.demands {
		if d.BuyerID == buyerID {
			result = append(result, d)
		}
	}
	return result, nil
}
func (f *fakeDemandRepo) Update(ctx context.Context, d *Demand) error {
	f.demands[d.ID] = *d
	return nil
}
func (f *fakeDemandRepo) Cancel(ctx context.Context, id uint) error {
	d := f.demands[id]
	d.Status = "CANCELLED"
	f.demands[id] = d
	return nil
}
func (f *fakeDemandRepo) List(ctx context.Context, filters DemandFilters) ([]Demand, error) {
	return nil, nil
}

func newFakeRepo() *fakeDemandRepo {
	return &fakeDemandRepo{demands: make(map[uint]Demand)}
}

func TestCreateDemand_Success(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo)
	ctx := context.Background()

	d := &Demand{
		CropName:    "Wheat",
		Quantity:    100,
		Unit:        "kg",
		TargetPrice: 2000,
		State:       "UP",
		District:    "Agra",
		RequiredBy:  time.Now().AddDate(0, 0, 7),
	}

	err := svc.CreateDemand(ctx, 10, d)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	stored, _ := repo.FindByID(ctx, d.ID)
	if stored.BuyerID != 10 {
		t.Errorf("expected BuyerID 10, got %d", stored.BuyerID)
	}
	if stored.Status != "ACTIVE" {
		t.Errorf("expected status ACTIVE, got %s", stored.Status)
	}
}

func TestCreateDemand_ValidationCropNameEmpty(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo)
	ctx := context.Background()

	d := &Demand{Quantity: 100, Unit: "kg", TargetPrice: 2000, State: "UP", District: "Agra", RequiredBy: time.Now().AddDate(0, 0, 7)}
	err := svc.CreateDemand(ctx, 10, d)
	if err != ErrInvalidCropName {
		t.Errorf("expected ErrInvalidCropName, got %v", err)
	}
}

func TestCreateDemand_ValidationQuantityZero(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo)
	ctx := context.Background()

	d := &Demand{CropName: "Wheat", Quantity: 0, Unit: "kg", TargetPrice: 2000, State: "UP", District: "Agra", RequiredBy: time.Now().AddDate(0, 0, 7)}
	err := svc.CreateDemand(ctx, 10, d)
	if err != ErrInvalidQuantity {
		t.Errorf("expected ErrInvalidQuantity, got %v", err)
	}
}

func TestCreateDemand_Unauthorized(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo)
	ctx := context.Background()

	d := &Demand{CropName: "Wheat", Quantity: 100, Unit: "kg", TargetPrice: 2000, State: "UP", District: "Agra", RequiredBy: time.Now().AddDate(0, 0, 7)}
	err := svc.CreateDemand(ctx, 0, d)
	if err != ErrUnauthorized {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}

func TestGetDemand_Success(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo)
	ctx := context.Background()

	d := &Demand{ID: 1, BuyerID: 10, CropName: "Wheat", Quantity: 100, Unit: "kg", TargetPrice: 2000, State: "UP", District: "Agra", RequiredBy: time.Now(), Status: "ACTIVE"}
	repo.Create(ctx, d)

	fetched, err := svc.GetDemand(ctx, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if fetched.ID != 1 {
		t.Errorf("expected ID 1, got %d", fetched.ID)
	}
}

func TestGetDemand_NotFound(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo)
	ctx := context.Background()

	_, err := svc.GetDemand(ctx, 999)
	if err != ErrDemandNotFound {
		t.Errorf("expected ErrDemandNotFound, got %v", err)
	}
}

func TestGetBuyerDemands_Success(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo)
	ctx := context.Background()

	d1 := &Demand{ID: 1, BuyerID: 10, CropName: "Wheat", Status: "ACTIVE", RequiredBy: time.Now()}
	d2 := &Demand{ID: 2, BuyerID: 10, CropName: "Onion", Status: "ACTIVE", RequiredBy: time.Now()}
	d3 := &Demand{ID: 3, BuyerID: 20, CropName: "Rice", Status: "ACTIVE", RequiredBy: time.Now()}
	repo.Create(ctx, d1)
	repo.Create(ctx, d2)
	repo.Create(ctx, d3)

	demands, err := svc.GetBuyerDemands(ctx, 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(demands) != 2 {
		t.Errorf("expected 2 demands for buyer 10, got %d", len(demands))
	}
}

func TestGetBuyerDemands_Unauthorized(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo)
	ctx := context.Background()

	_, err := svc.GetBuyerDemands(ctx, 0)
	if err != ErrUnauthorized {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}

func TestUpdateDemand_Success(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo)
	ctx := context.Background()

	d := &Demand{ID: 1, BuyerID: 10, CropName: "Wheat", Quantity: 100, Unit: "kg", TargetPrice: 2000, State: "UP", District: "Agra", Status: "ACTIVE", RequiredBy: time.Now().AddDate(0, 0, 7)}
	repo.Create(ctx, d)

	update := &Demand{ID: 1, CropName: "Wheat", Quantity: 150, Unit: "kg", TargetPrice: 2100, State: "UP", District: "Agra", RequiredBy: time.Now().AddDate(0, 0, 7)}
	err := svc.UpdateDemand(ctx, 10, update)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updated, _ := repo.FindByID(ctx, 1)
	if updated.Quantity != 150 {
		t.Errorf("expected quantity 150, got %f", updated.Quantity)
	}
}

func TestUpdateDemand_Unauthorized(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo)
	ctx := context.Background()

	d := &Demand{ID: 1, BuyerID: 10, CropName: "Wheat", Quantity: 100, Unit: "kg", TargetPrice: 2000, State: "UP", District: "Agra", Status: "ACTIVE", RequiredBy: time.Now().AddDate(0, 0, 7)}
	repo.Create(ctx, d)

	update := &Demand{ID: 1, CropName: "Wheat", Quantity: 150, Unit: "kg", TargetPrice: 2100, State: "UP", District: "Agra", RequiredBy: time.Now().AddDate(0, 0, 7)}
	err := svc.UpdateDemand(ctx, 99, update)
	if err != ErrUnauthorized {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}

func TestUpdateDemand_InactiveStatus(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo)
	ctx := context.Background()

	d := &Demand{ID: 1, BuyerID: 10, CropName: "Wheat", Quantity: 100, Unit: "kg", TargetPrice: 2000, State: "UP", District: "Agra", Status: "COMPLETED", RequiredBy: time.Now().AddDate(0, 0, 7)}
	repo.Create(ctx, d)

	update := &Demand{ID: 1, CropName: "Wheat", Quantity: 150, Unit: "kg", TargetPrice: 2100, State: "UP", District: "Agra", RequiredBy: time.Now().AddDate(0, 0, 7)}
	err := svc.UpdateDemand(ctx, 10, update)
	if err != ErrDemandNotActive {
		t.Errorf("expected ErrDemandNotActive, got %v", err)
	}
}

func TestCancelDemand_Success(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo)
	ctx := context.Background()

	d := &Demand{ID: 1, BuyerID: 10, CropName: "Wheat", Status: "ACTIVE", RequiredBy: time.Now()}
	repo.Create(ctx, d)

	err := svc.CancelDemand(ctx, 10, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	cancelled, _ := repo.FindByID(ctx, 1)
	if cancelled.Status != "CANCELLED" {
		t.Errorf("expected status CANCELLED, got %s", cancelled.Status)
	}
}

func TestCancelDemand_Unauthorized(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo)
	ctx := context.Background()

	d := &Demand{ID: 1, BuyerID: 10, CropName: "Wheat", Status: "ACTIVE", RequiredBy: time.Now()}
	repo.Create(ctx, d)

	err := svc.CancelDemand(ctx, 99, 1)
	if err != ErrUnauthorized {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}

func TestCancelDemand_InactiveStatus(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo)
	ctx := context.Background()

	d := &Demand{ID: 1, BuyerID: 10, CropName: "Wheat", Status: "COMPLETED", RequiredBy: time.Now()}
	repo.Create(ctx, d)

	err := svc.CancelDemand(ctx, 10, 1)
	if err != ErrDemandNotActive {
		t.Errorf("expected ErrDemandNotActive, got %v", err)
	}
}

func TestCancelDemand_NotFound(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo)
	ctx := context.Background()

	err := svc.CancelDemand(ctx, 10, 999)
	if err != ErrDemandNotFound {
		t.Errorf("expected ErrDemandNotFound, got %v", err)
	}
}
