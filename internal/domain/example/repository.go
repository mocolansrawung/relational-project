package example

//go:generate go run github.com/golang/mock/mockgen -source repository.go -destination mock/repository_mock.go -package example_mock

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

// SomeRepository is the repository for SomeEntity data
type SomeRepository interface {
	ResolveByID(id uuid.UUID) (someEntity SomeEntity, err error)
}

// SomeRepositoryMySQL is the MySQL-backed implementation of SomeRepository
type SomeRepositoryMySQL struct {
	DB *infras.MySQLConn `inject:"db"`
}

// ResolveByID resolves a SomeEntity by its ID
func (r *SomeRepositoryMySQL) ResolveByID(id uuid.UUID) (someEntity SomeEntity, err error) {
	switch id.String() {
	case "00000000-0000-0000-0000-000000000000":
		// Test case demonstrating a NOT FOUND scenario
		return someEntity, failure.NotFound("someEntity")
	case "11111111-1111-1111-1111-111111111111":
		// Test case demonstrating an INTERNAL SERVER ERROR scenario
		return someEntity, errors.New("some random internal server error")
	default:
		err = r.DB.Read.Get(&someEntity, fmt.Sprintf(querySample, id.String()))
	}
	return
}
