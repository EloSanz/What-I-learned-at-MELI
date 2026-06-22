package item

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var (
	ErrItemNotFound = errors.New("item not found")
	ErrOutOfStock   = errors.New("insufficient stock available")
)

type Repository interface {
	FindByID(id string) (*Item, error)
	DecrementStock(id string, quantity int) (*Item, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new instance of the items repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// FindByID retrieves an item by its primary ID
func (r *repository) FindByID(id string) (*Item, error) {
	var item Item
	err := r.db.First(&item, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrItemNotFound
		}
		return nil, err
	}
	return &item, nil
}

// DecrementStock reduces the stock of an item in the database within a transaction
func (r *repository) DecrementStock(id string, quantity int) (*Item, error) {
	var item Item
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Find the item and lock the row for writing (Select for Update)
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&item, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrItemNotFound
			}
			return err
		}

		// Validate stock
		if item.Stock < quantity {
			return ErrOutOfStock
		}

		// Decrement stock
		item.Stock -= quantity
		if err := tx.Save(&item).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

	return &item, nil
}
