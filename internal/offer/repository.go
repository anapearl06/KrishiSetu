package offer

import "gorm.io/gorm"

type Repository interface {
	Create(offer *Offer) error
	FindByID(id uint) (*Offer, error)
	FindByBuyer(buyerID uint) ([]Offer, error)
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

func (r *repository) FindByBuyer(buyerID uint) ([]Offer, error) {
	var offers []Offer

	err := r.db.
		Where("buyer_id = ?", buyerID).
		Order("created_at DESC").
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