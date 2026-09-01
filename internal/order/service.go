package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/raaj2493/KrishiSetu/internal/listing"
	"github.com/raaj2493/KrishiSetu/internal/offer"
	"gorm.io/gorm"
)

var (
	ErrOrderNotFound       = errors.New("order not found")
	ErrOfferNotFound       = errors.New("offer not found")
	ErrOfferNotPending     = errors.New("offer is not pending")
	ErrListingNotFound     = errors.New("listing not found")
	ErrListingNotActive    = errors.New("listing is not active")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrQuantityUnavailable = errors.New("quantity no longer available")
)

type Service interface {
	AcceptOffer(c context.Context, offerID uint, farmerID uint) (*Order, error)

	GetOrder(id uint) (*Order, error)

	GetBuyerOrders(buyerID uint) ([]OrderView, error)

	GetFarmerOrders(farmerID uint) ([]OrderView, error)
}

type service struct {
	repo        Repository
	offerRepo   offer.Repository
	listingRepo listing.Repository
	db          *gorm.DB
}

func NewService(
	repo Repository,
	offerRepo offer.Repository,
	listingRepo listing.Repository,
	db *gorm.DB,
) Service {
	return &service{
		repo:        repo,
		offerRepo:   offerRepo,
		listingRepo: listingRepo,
		db:          db,
	}
}

func (s *service) AcceptOffer(
	c context.Context,
	offerID uint,
	farmerID uint,
) (*Order, error) {

	var createdOrder *Order

	err := s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {

		// --------------------------------
		// 1. Get Offer
		// --------------------------------

		var offerItem offer.Offer

		if err := tx.WithContext(c).First(&offerItem, offerID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrOfferNotFound
			}

			return fmt.Errorf("find offer: %w", err)
		}

		// Offer must still be pending.
		if offerItem.Status != "PENDING" {
			return ErrOfferNotPending
		}

		// --------------------------------
		// 2. Get Listing
		// --------------------------------

		var listingItem listing.CropListing

		if err := tx.WithContext(c).First(
			&listingItem,
			offerItem.ListingID,
		).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrListingNotFound
			}

			return fmt.Errorf("find listing: %w", err)
		}

		// --------------------------------
		// 3. Verify Farmer Ownership
		// --------------------------------

		if listingItem.FarmerID != farmerID {
			return ErrUnauthorized
		}

		// --------------------------------
		// 4. Verify Listing
		// --------------------------------

		if listingItem.Status != "ACTIVE" {
			return ErrListingNotActive
		}

		if offerItem.Quantity > listingItem.Quantity {
			return ErrQuantityUnavailable
		}

		// --------------------------------
		// 5. Create Order
		// --------------------------------

		order := &Order{
			OfferID:     offerItem.ID,
			ListingID:   listingItem.ID,
			BuyerID:     offerItem.BuyerID,
			FarmerID:    listingItem.FarmerID,
			Quantity:    offerItem.Quantity,
			AgreedPrice: offerItem.OfferedPrice,
			TotalAmount: offerItem.Quantity * offerItem.OfferedPrice,
			Status:      "CONFIRMED",
		}

		if err := tx.WithContext(c).Create(order).Error; err != nil {
			return fmt.Errorf("create order: %w", err)
		}

		// --------------------------------
		// 6. Reduce Listing Quantity
		// --------------------------------

		listingItem.Quantity -= offerItem.Quantity

		// If nothing remains, mark listing sold.
		if listingItem.Quantity == 0 {
			listingItem.Status = "SOLD"
		}

		if err := tx.WithContext(c).Save(&listingItem).Error; err != nil {
			return fmt.Errorf("update listing: %w", err)
		}

		// --------------------------------
		// 7. Accept Offer
		// --------------------------------

		offerItem.Status = "ACCEPTED"

		if err := tx.WithContext(c).Save(&offerItem).Error; err != nil {
			return fmt.Errorf("update offer: %w", err)
		}

		createdOrder = order

		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdOrder, nil
}

func (s *service) GetOrder(id uint) (*Order, error) {
	order, err := s.repo.FindByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}

		return nil, fmt.Errorf("find order: %w", err)
	}

	return order, nil
}

func (s *service) GetBuyerOrders(buyerID uint) ([]OrderView, error) {
	return s.repo.FindByBuyer(buyerID)
}

func (s *service) GetFarmerOrders(farmerID uint) ([]OrderView, error) {
	return s.repo.FindByFarmer(farmerID)
}
