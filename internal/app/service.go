package app

type Service struct {
	storage Facade
}

func NewService(storage Facade) *Service {
	return &Service{storage: storage}
}
