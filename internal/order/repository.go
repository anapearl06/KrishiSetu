package order

import "gorm.io/gorm"

type Repository interface {
	Create(order *Order) error
	FindByID(id uint) (*Order, error)
	FindByBuyer(buyerID uint) ([]Order, error)
	FindByFarmer(farmerID uint) ([]Order, error)
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

func (r *repository) FindByBuyer(buyerID uint) ([]Order, error) {
	var orders []Order

	err := r.db.
		Where("buyer_id = ?", buyerID).
		Order("created_at DESC").
		Find(&orders).Error

	return orders, err
}

func (r *repository) FindByFarmer(farmerID uint) ([]Order, error) {
	var orders []Order

	err := r.db.
		Where("farmer_id = ?", farmerID).
		Order("created_at DESC").
		Find(&orders).Error

	return orders, err
}

func (r *repository) Update(order *Order) error {
	return r.db.Save(order).Error
}