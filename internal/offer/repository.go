package offer

import "gorm.io/gorm"

type Repository interface {
	Create(o *Offer) error
	FindByID(id uint) (*Offer, error)
	FindByBuyer(buyerID uint) ([]Offer, error)
	FindByFarmer(farmerID uint) ([]Offer, error)
	FindByListing(listingID uint) ([]Offer, error)
	Update(o *Offer) error
}

type postgresRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(o *Offer) error {
	return r.db.Create(o).Error
}

func (r *postgresRepository) FindByID(id uint) (*Offer, error) {
	var o Offer
	if err := r.db.First(&o, id).Error; err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *postgresRepository) FindByBuyer(buyerID uint) ([]Offer, error) {
	var offers []Offer
	if err := r.db.Where("buyer_id = ?", buyerID).Order("created_at DESC").Find(&offers).Error; err != nil {
		return nil, err
	}
	return offers, nil
}

func (r *postgresRepository) FindByFarmer(farmerID uint) ([]Offer, error) {
	var offers []Offer
	if err := r.db.Where("farmer_id = ?", farmerID).Order("created_at DESC").Find(&offers).Error; err != nil {
		return nil, err
	}
	return offers, nil
}

func (r *postgresRepository) FindByListing(listingID uint) ([]Offer, error) {
	var offers []Offer
	if err := r.db.Where("listing_id = ?", listingID).Order("created_at DESC").Find(&offers).Error; err != nil {
		return nil, err
	}
	return offers, nil
}

func (r *postgresRepository) Update(o *Offer) error {
	return r.db.Save(o).Error
}
