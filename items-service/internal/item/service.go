package item

type Service interface {
	GetByID(id string) (*Item, error)
	ValidateAndReserveStock(id string, quantity int) (*Item, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetByID(id string) (*Item, error) {
	return s.repo.FindByID(id)
}

func (s *service) ValidateAndReserveStock(id string, quantity int) (*Item, error) {
	return s.repo.DecrementStock(id, quantity)
}
