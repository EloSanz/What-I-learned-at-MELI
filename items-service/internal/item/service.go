package item

import (
	"items-service/internal/infra"
)

type Service interface {
	GetByID(id string) (*Item, error)
	ValidateAndReserveStock(id string, quantity int) (*Item, error)
}

type service struct {
	repo        Repository
	lockService infra.LockService
}

func NewService(repo Repository, lockService infra.LockService) Service {
	return &service{
		repo:        repo,
		lockService: lockService,
	}
}

func (s *service) GetByID(id string) (*Item, error) {
	return s.repo.FindByID(id)
}

func (s *service) ValidateAndReserveStock(id string, quantity int) (*Item, error) {
	var item *Item
	err := s.lockService.WithLock(id, func() error {
		var innerErr error
		item, innerErr = s.repo.DecrementStock(id, quantity)
		return innerErr
	})
	return item, err
}
