package offer

import "gorm.io/gorm"

type Repository interface {
	Create(offer *Offer) error
	FindByID(id uint) (*Offer, error)
	FindByBuyer(buyerID uint) ([]OfferView, error)
	FindByFarmer(farmerID uint) ([]OfferView, error)
	FindByListing(listingID uint) ([]Offer, error)
	Update(offer *Offer) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(offer *Offer) error {
	return r.db.Create(offer).Error
}

func (r *repository) FindByID(id uint) (*Offer, error) {
	var offer Offer

	err := r.db.First(&offer, id).Error
	if err != nil {
		return nil, err
	}

	return &offer, nil
}

func (r *repository) FindByBuyer(buyerID uint) ([]OfferView, error) {
	var offers []OfferView

	err := r.db.
		Table("offers").
		Select("offers.*, crop_listings.crop_name AS crop").
		Joins("JOIN crop_listings ON crop_listings.id = offers.listing_id").
		Where("offers.buyer_id = ?", buyerID).
		Order("offers.created_at DESC").
		Find(&offers).Error

	return offers, err
}

func (r *repository) FindByFarmer(farmerID uint) ([]OfferView, error) {
	var offers []OfferView

	err := r.db.
		Table("offers").
		Select("offers.*, crop_listings.crop_name AS crop, buyers.name AS buyer_name").
		Joins("JOIN crop_listings ON crop_listings.id = offers.listing_id").
		Joins("JOIN buyers ON buyers.id = offers.buyer_id").
		Where("crop_listings.farmer_id = ?", farmerID).
		Order("offers.created_at DESC").
		Find(&offers).Error

	return offers, err
}

func (r *repository) FindByListing(listingID uint) ([]Offer, error) {
	var offers []Offer

	err := r.db.
		Where("listing_id = ?", listingID).
		Order("created_at DESC").
		Find(&offers).Error

	return offers, err
}

func (r *repository) Update(offer *Offer) error {
	return r.db.Save(offer).Error
}