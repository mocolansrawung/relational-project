package foobarbaz

//go:generate go run github.com/golang/mock/mockgen -source foo_service.go -destination mock/foo_service_mock.go -package foobarbaz_mock

import (
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/gofrs/uuid"
)

// FooService is the service interface for Foo entities.
type FooService interface {
	Create(requestFormat FooRequestFormat, userID uuid.UUID) (foo Foo, err error)
	ResolveByID(id uuid.UUID, withItems bool) (foo Foo, err error)
}

// FooServiceImpl is the service implementation for Foo entities.
type FooServiceImpl struct {
	FooRepository FooRepository
}

// ProvideFooServiceImpl is the provider for this service.
func ProvideFooServiceImpl(fooRepository FooRepository) *FooServiceImpl {
	s := new(FooServiceImpl)
	s.FooRepository = fooRepository
	return s
}

// Create creates a new Foo.
func (s *FooServiceImpl) Create(requestFormat FooRequestFormat, userID uuid.UUID) (foo Foo, err error) {
	foo = foo.NewFromRequestFormat(requestFormat, userID)

	foo.Recalculate()
	err = foo.Validate()
	if err != nil {
		return foo, failure.BadRequest(err)
	}

	err = s.FooRepository.Create(foo)
	return
}

// ResolveByID resolves a Foo by its ID.
func (s *FooServiceImpl) ResolveByID(id uuid.UUID, withItems bool) (foo Foo, err error) {
	foo, err = s.FooRepository.ResolveByID(id)

	if withItems {
		items, err := s.FooRepository.ResolveItemsByFooIDs([]uuid.UUID{foo.ID})
		if err != nil {
			return foo, err
		}

		foo.AttachItems(items)
	}

	return
}
