package order

import (
	"crypto/rand"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Possible Order statuses
const (
	StatusPending        = "PENDING"
	StatusReadyToProcess = "READY_TO_PROCESS"
	StatusCompleted      = "COMPLETED"
	StatusFailed         = "FAILED"
)

// Order represents the model of a purchase in the system
type Order struct {
	ID        string         `gorm:"primaryKey;type:varchar(50)" json:"id"`
	UserID    string         `gorm:"type:varchar(100);not null" json:"user_id"`
	ItemID    string         `gorm:"type:varchar(50);not null" json:"item_id"`
	Quantity  int            `gorm:"type:integer;not null" json:"quantity"`
	Amount    float64        `gorm:"type:numeric(12,2);not null" json:"amount"`
	Address   string         `gorm:"type:varchar(255);not null" json:"address"`
	Status    string         `gorm:"type:varchar(30);not null;default:'PENDING'" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// GenerateUUID generates a random unique identifier (basic v4) without dependencies
func GenerateUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
