package infra

import (
	"errors"
	"fmt"
	"hash/fnv"

	"gorm.io/gorm"
)

// ErrResourceLocked is returned when an advisory lock cannot be acquired
var ErrResourceLocked = errors.New("resource is currently locked")

// LockService defines the interface for acquiring distributed locks
type LockService interface {
	// WithLock executes the given function only if the lock for resourceID can be acquired.
	// It uses a transaction to ensure the lock is automatically released upon completion.
	WithLock(resourceID string, fn func() error) error
}

type pgLockService struct {
	db *gorm.DB
}

// NewPGLockService creates a new PostgreSQL advisory lock service
func NewPGLockService(db *gorm.DB) LockService {
	return &pgLockService{db: db}
}

func (s *pgLockService) WithLock(resourceID string, fn func() error) error {
	// Hash the resource ID to an int64 for the advisory lock
	h := fnv.New64a()
	h.Write([]byte(resourceID))
	lockID := int64(h.Sum64())

	return s.db.Transaction(func(tx *gorm.DB) error {
		var locked bool
		// pg_try_advisory_xact_lock acquires a transaction-level lock.
		// It returns true if successful, false otherwise.
		// The lock is automatically released at the end of the transaction.
		err := tx.Raw("SELECT pg_try_advisory_xact_lock(?)", lockID).Scan(&locked).Error
		if err != nil {
			return fmt.Errorf("failed to acquire advisory lock: %w", err)
		}

		if !locked {
			return ErrResourceLocked
		}

		// Lock acquired successfully, execute the critical section
		return fn()
	})
}
