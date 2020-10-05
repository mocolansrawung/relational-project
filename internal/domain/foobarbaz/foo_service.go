package foobarbaz

//go:generate go run github.com/golang/mock/mockgen -source foo_service.go -destination mock/foo_service_mock.go -package foobarbaz_mock

import (
	"github.com/gofrs/uuid"
)

type FooService interface {
	ResolveByID(id uuid.UUID) (Foo, error)
}

type FooServiceImpl struct {
	FooRepository FooRepository
}

func ProvideFooServiceImpl(fooRepository FooRepository) *FooServiceImpl {
	s := new(FooServiceImpl)
	s.FooRepository = fooRepository
	return s
}

func (s *FooServiceImpl) ResolveByID(id uuid.UUID) (Foo, error) {
	return s.FooRepository.ResolveByID(id)
}
