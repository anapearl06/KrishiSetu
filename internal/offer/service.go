package offer

import "errors"

type Service struct {
	repo Repository
}

type CreateInput struct {
	ListingID uint    `json:"listing_id"`
	Price     float64 `json:"price"`
	Quantity  float64 `json:"quantity"`
	Message   string  `json:"message"`
}

type RespondInput struct {
	Action string `json:"action"`
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateOffer(buyerID, farmerID uint, input CreateInput) (*Offer, error) {
	if input.ListingID == 0 {
		return nil, errors.New("listing_id is required")
	}
	if input.Price <= 0 {
		return nil, errors.New("price must be greater than zero")
	}
	if input.Quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}

	offer := &Offer{
		ListingID: input.ListingID,
		FarmerID:  farmerID,
		BuyerID:   buyerID,
		Price:     input.Price,
		Quantity:  input.Quantity,
		Message:   input.Message,
		Status:    "PENDING",
	}

	if err := s.repo.Create(offer); err != nil {
		return nil, err
	}

	return offer, nil
}

func (s *Service) GetBuyerOffers(buyerID uint) ([]Offer, error) {
	return s.repo.FindByBuyer(buyerID)
}

func (s *Service) GetFarmerOffers(farmerID uint) ([]Offer, error) {
	return s.repo.FindByFarmer(farmerID)
}

func (s *Service) GetListingOffers(listingID uint) ([]Offer, error) {
	return s.repo.FindByListing(listingID)
}

func (s *Service) RespondOffer(farmerID, offerID uint, input RespondInput) (*Offer, error) {
	offer, err := s.repo.FindByID(offerID)
	if err != nil {
		return nil, errors.New("offer not found")
	}

	if offer.FarmerID != farmerID {
		return nil, errors.New("unauthorized: only the listing owner can respond")
	}

	switch input.Action {
	case "ACCEPT":
		offer.Status = "ACCEPTED"
	case "REJECT":
		offer.Status = "REJECTED"
	default:
		return nil, errors.New("action must be ACCEPT or REJECT")
	}

	if err := s.repo.Update(offer); err != nil {
		return nil, err
	}

	return offer, nil
}
