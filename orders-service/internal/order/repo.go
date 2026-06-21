package order

import (
	"errors"

	"gorm.io/gorm"
)

var ErrOrderNotFound = errors.New("order not found")

type Repository interface {
	Create(order *Order) error
	FindByID(id string) (*Order, error)
	UpdateStatus(id string, status string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create inserta una nueva orden en la base de datos
func (r *repository) Create(order *Order) error {
	return r.db.Create(order).Error
}

// FindByID busca una orden por su identificador primario
func (r *repository) FindByID(id string) (*Order, error) {
	var ord Order
	err := r.db.First(&ord, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}
	return &ord, nil
}

// UpdateStatus actualiza el estado de una orden específica
func (r *repository) UpdateStatus(id string, status string) error {
	result := r.db.Model(&Order{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrOrderNotFound
	}
	return nil
}
