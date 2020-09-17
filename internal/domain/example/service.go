package example

//go:generate mockgen -source service.go -destination mock/service_mock.go -package example_mock

import "github.com/gofrs/uuid"

type SomeService interface {
	ResolveByID(id uuid.UUID) (SomeEntity, error)
}

type SomeServiceImpl struct {
	SomeRepository SomeRepository `inject:"example.someRepository"`
}

func (s *SomeServiceImpl) ResolveByID(id uuid.UUID) (SomeEntity, error) {
	return s.SomeRepository.ResolveByID(id)
}
