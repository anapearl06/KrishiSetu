package order

import "gorm.io/gorm"

type Repository interface {
	Create(order *Order) error
	FindByID(id uint) (*Order, error)
	FindByBuyer(buyerID uint) ([]OrderView, error)
	FindByFarmer(farmerID uint) ([]OrderView, error)
	Update(order *Order) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(order *Order) error {
	return r.db.Create(order).Error
}

func (r *repository) FindByID(id uint) (*Order, error) {
	var order Order

	err := r.db.First(&order, id).Error
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *repository) FindByBuyer(buyerID uint) ([]OrderView, error) {
	var orders []OrderView

	err := r.db.
		Table("orders").
		Select("orders.*, crop_listings.crop_name AS crop, farmers.name AS farmer_name").
		Joins("JOIN crop_listings ON crop_listings.id = orders.listing_id").
		Joins("JOIN farmers ON farmers.id = orders.farmer_id").
		Where("orders.buyer_id = ?", buyerID).
		Order("orders.created_at DESC").
		Find(&orders).Error

	return orders, err
}

func (r *repository) FindByFarmer(farmerID uint) ([]OrderView, error) {
	var orders []OrderView

	err := r.db.
		Table("orders").
		Select("orders.*, crop_listings.crop_name AS crop, farmers.name AS farmer_name").
		Joins("JOIN crop_listings ON crop_listings.id = orders.listing_id").
		Joins("JOIN farmers ON farmers.id = orders.farmer_id").
		Where("orders.farmer_id = ?", farmerID).
		Order("orders.created_at DESC").
		Find(&orders).Error

	return orders, err
}

func (r *repository) Update(order *Order) error {
	return r.db.Save(order).Error
}