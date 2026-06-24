package item

import (
	"time"

	"gorm.io/gorm"
)

// Item represents the domain model of a Product/Item in Mercado Libre
type Item struct {
	ID        string         `gorm:"primaryKey;type:varchar(50)" json:"id"`
	Title     string         `gorm:"type:varchar(255);not null" json:"title"`
	Price     float64        `gorm:"type:numeric(12,2);not null" json:"price"`
	Stock     int            `gorm:"type:integer;not null;default:0" json:"stock"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type ValidateStockRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}
