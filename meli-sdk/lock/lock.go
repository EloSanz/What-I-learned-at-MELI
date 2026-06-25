package lock

import (
	"errors"
	"fmt"
	"hash/fnv"

	"gorm.io/gorm"
)

// ErrResourceLocked is returned when an advisory lock cannot be acquired
var ErrResourceLocked = errors.New("resource is currently locked")

// Service defines the interface for acquiring distributed locks
type Service interface {
	WithLock(resourceID string, fn func() error) error
}

type pgLockService struct {
	db *gorm.DB
}

// NewPGLockService creates a new PostgreSQL advisory lock service
func NewPGLockService(db *gorm.DB) Service {
	return &pgLockService{db: db}
}

func (s *pgLockService) WithLock(resourceID string, fn func() error) error {
	h := fnv.New64a()
	h.Write([]byte(resourceID))
	lockID := int64(h.Sum64())

	return s.db.Transaction(func(tx *gorm.DB) error {
		var locked bool
		err := tx.Raw("SELECT pg_try_advisory_xact_lock(?)", lockID).Scan(&locked).Error
		if err != nil {
			return fmt.Errorf("failed to acquire advisory lock: %w", err)
		}

		if !locked {
			return ErrResourceLocked
		}

		return fn()
	})
}
