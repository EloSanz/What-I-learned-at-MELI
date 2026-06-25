package item

import (
	"github.com/user/meli-sdk/lock"
)

type Service interface {
	GetByID(id string) (*Item, error)
	ValidateAndReserveStock(id string, quantity int) (*Item, error)
}

type service struct {
	repo        Repository
	lockService lock.Service
}

func NewService(repo Repository, lockService lock.Service) Service {
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
