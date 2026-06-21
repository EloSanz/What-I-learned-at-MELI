package order

import (
	"crypto/rand"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Estados posibles de una Orden
const (
	StatusPending        = "PENDING"
	StatusReadyToProcess = "READY_TO_PROCESS"
	StatusCompleted      = "COMPLETED"
	StatusFailed         = "FAILED"
)

// Order representa el modelo de una compra en el sistema
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

// GenerateUUID genera un identificador único aleatorio (versión básica v4) sin dependencias
func GenerateUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
