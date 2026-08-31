package offer

<<<<<<< HEAD
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
=======
import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var (
	ErrOfferNotFound       = errors.New("offer not found")
	ErrListingNotFound     = errors.New("listing not found")
	ErrListingNotActive    = errors.New("listing is not active")
	ErrInvalidQuantity     = errors.New("invalid quantity")
	ErrInvalidPrice        = errors.New("invalid price")
	ErrQuantityUnavailable = errors.New("requested quantity exceeds available quantity")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrInvalidStatus       = errors.New("invalid offer status")
)

type Service interface {
	CreateOffer(
		buyerID uint,
		listingID uint,
		quantity float64,
		offeredPrice float64,
		message string,
	) (*Offer, error)

	GetOffer(id uint) (*Offer, error)

	GetBuyerOffers(buyerID uint) ([]Offer, error)

	GetListingOffers(listingID uint) ([]Offer, error)

	CancelOffer(offerID uint, buyerID uint) error
}

type service struct {
	repo Repository
	db   *gorm.DB
}

func NewService(repo Repository, db *gorm.DB) Service {
	return &service{
		repo: repo,
		db:   db,
	}
}

func (s *service) CreateOffer(
	buyerID uint,
	listingID uint,
	quantity float64,
	offeredPrice float64,
	message string,
) (*Offer, error) {

	if quantity <= 0 {
		return nil, ErrInvalidQuantity
	}

	if offeredPrice < 0 {
		return nil, errors.New("offered price cannot be negative")
	}

	// We need the listing to validate availability.
	var listing struct {
		ID       uint
		Quantity float64
		Status   string
	}

	err := s.db.
		Table("crop_listings").
		Select("id, quantity, status").
		Where("id = ?", listingID).
		First(&listing).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrListingNotFound
		}

		return nil, fmt.Errorf("find listing: %w", err)
	}

	if listing.Status != "ACTIVE" {
		return nil, ErrListingNotActive
	}

	if quantity > listing.Quantity {
		return nil, ErrQuantityUnavailable
	}

	offer := &Offer{
		ListingID:    listingID,
		BuyerID:      buyerID,
		Quantity:     quantity,
		OfferedPrice: offeredPrice,
		Message:      message,
		Status:       "PENDING",
	}

	if err := s.repo.Create(offer); err != nil {
		return nil, fmt.Errorf("create offer: %w", err)
>>>>>>> tmp-pr-merge
	}

	return offer, nil
}

<<<<<<< HEAD
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
=======
func (s *service) GetOffer(id uint) (*Offer, error) {
	offer, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOfferNotFound
		}

		return nil, fmt.Errorf("find offer: %w", err)
	}

	return offer, nil
}

func (s *service) GetBuyerOffers(buyerID uint) ([]Offer, error) {
	return s.repo.FindByBuyer(buyerID)
}

func (s *service) GetListingOffers(listingID uint) ([]Offer, error) {
	return s.repo.FindByListing(listingID)
}

func (s *service) CancelOffer(offerID uint, buyerID uint) error {
	offer, err := s.repo.FindByID(offerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrOfferNotFound
		}

		return fmt.Errorf("find offer: %w", err)
	}

	if offer.BuyerID != buyerID {
		return ErrUnauthorized
	}

	if offer.Status != "PENDING" {
		return ErrInvalidStatus
	}

	offer.Status = "CANCELLED"

	if err := s.repo.Update(offer); err != nil {
		return fmt.Errorf("cancel offer: %w", err)
	}

	return nil
>>>>>>> tmp-pr-merge
}
