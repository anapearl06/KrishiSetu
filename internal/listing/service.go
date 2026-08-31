package listing

import (
	"errors"
	"time"
)

type Service struct {
	repo Repository
}

type CreateInput struct {
	Crop        string  `json:"crop"`
	Quantity    float64 `json:"quantity"`
	Unit        string  `json:"unit"`
	Price       float64 `json:"price"`
	State       string  `json:"state"`
	District    string  `json:"district"`
	Description string  `json:"description"`
}

type UpdateInput struct {
	Price    float64 `json:"price"`
	Quantity float64 `json:"quantity"`
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateListing(farmerID uint, input CreateInput) (*CropListing, error) {
	if input.Crop == "" {
		return nil, errors.New("crop name is required")
	}
	if input.Quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}
	if input.Unit == "" {
		return nil, errors.New("unit is required")
	}
	if input.Price <= 0 {
		return nil, errors.New("price must be greater than zero")
	}
	if input.State == "" {
		return nil, errors.New("state is required")
	}
	if input.District == "" {
		return nil, errors.New("district is required")
	}

	listing := &CropListing{
		FarmerID:      farmerID,
		CropName:      input.Crop,
		Quantity:      input.Quantity,
		Unit:          input.Unit,
		ExpectedPrice: input.Price,
		Description:   input.Description,
		State:         input.State,
		District:      input.District,
		HarvestDate:   time.Now(),
		Status:        "ACTIVE",
	}

	if err := s.repo.Create(listing); err != nil {
		return nil, err
	}

	return listing, nil
}

func (s *Service) GetMyListings(farmerID uint) ([]CropListing, error) {
	return s.repo.FindByFarmerID(farmerID)
}

func (s *Service) BrowseListings(crop, state, status string) ([]CropListing, error) {
	return s.repo.FindAll(crop, state, status)
}

func (s *Service) GetListing(id uint) (*CropListing, error) {
	return s.repo.FindByID(id)
}

func (s *Service) UpdateListing(id, farmerID uint, input UpdateInput) (*CropListing, error) {
	listing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if listing.FarmerID != farmerID {
		return nil, errors.New("unauthorized: you can only update your own listing")
	}

	if input.Price > 0 {
		listing.ExpectedPrice = input.Price
	}
	if input.Quantity > 0 {
		listing.Quantity = input.Quantity
	}

	if err := s.repo.Update(listing); err != nil {
		return nil, err
	}

	return listing, nil
}

func (s *Service) DeleteListing(id, farmerID uint) error {
	listing, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if listing.FarmerID != farmerID {
		return errors.New("unauthorized: you can only delete your own listing")
	}

	return s.repo.Delete(id)
}
