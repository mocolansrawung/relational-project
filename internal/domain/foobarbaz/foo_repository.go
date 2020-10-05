package foobarbaz

//go:generate go run github.com/golang/mock/mockgen -source foo_repository.go -destination mock/foo_repository_mock.go -package foobarbaz_mock

import (
	"errors"
	"fmt"

	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/gofrs/uuid"
)

const (
	querySample = `
		SELECT
			'%s' as id,
			'John Doe' as name,
			'active' as status,
			NOW() as created,
			'9c2a76af-0859-49f5-bf2f-06850227136a' as created_by,
			NOW() as updated,
			'0c2a76af-0859-49f5-bf2f-06850227136a' as updated_by`
)

// FooRepository is the repository for Foo data.
type FooRepository interface {
	ResolveByID(id uuid.UUID) (foo Foo, err error)
}

// FooRepositoryMySQL is the MySQL-backed implementation of FooRepository.
type FooRepositoryMySQL struct {
	DB *infras.MySQLConn
}

// ProvideFooRepositoryMySQL is the provider for this repository.
func ProvideFooRepositoryMySQL(db *infras.MySQLConn) *FooRepositoryMySQL {
	s := new(FooRepositoryMySQL)
	s.DB = db
	return s
}

// ResolveByID resolves a Foo by its ID
func (r *FooRepositoryMySQL) ResolveByID(id uuid.UUID) (foo Foo, err error) {
	switch id.String() {
	case "00000000-0000-0000-0000-000000000000":
		// Test case demonstrating a NOT FOUND scenario
		return foo, failure.NotFound("foo")
	case "11111111-1111-1111-1111-111111111111":
		// Test case demonstrating an INTERNAL SERVER ERROR scenario
		return foo, errors.New("foo internal server error")
	default:
		err = r.DB.Read.Get(&foo, fmt.Sprintf(querySample, id.String()))
	}
	return
}
