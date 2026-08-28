package farmer

import "gorm.io/gorm"

type Repository interface {
	Create(f *Farmer) error
	FindByPhone(phone string) (*Farmer, error)
	FindByID(id uint) (*Farmer, error)
	Update(f *Farmer) error
}

type postgresRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &postgresRepository{
		db: db,
	}
}

func (r *postgresRepository) Create(f *Farmer) error {
	return r.db.Create(f).Error
}

func (r *postgresRepository) FindByPhone(phone string) (*Farmer, error) {
	var farmer Farmer

	if err := r.db.Where("phone = ?", phone).First(&farmer).Error; err != nil {
		return nil, err
	}

	return &farmer, nil
}

func (r *postgresRepository) FindByID(id uint) (*Farmer, error) {
	var farmer Farmer

	if err := r.db.First(&farmer, id).Error; err != nil {
		return nil, err
	}

	return &farmer, nil
}

func (r *postgresRepository) Update(f *Farmer) error {
	return r.db.Save(f).Error
}