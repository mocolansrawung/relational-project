package example

//go:generate mockgen -source repository.go -destination mock/repository_mock.go -package example_mock

import (
	"fmt"

	"github.com/evermos/boilerplate-go/infras"
	"github.com/gofrs/uuid"
)

const (
	querySample = `
		SELECT
			'%s' as id,
			'John Doe' as name,
			NOW() as created,
			'9c2a76af-0859-49f5-bf2f-06850227136a' as created_by,
			NOW() as updated,
			'0c2a76af-0859-49f5-bf2f-06850227136a' as updated_by`
)

type SomeRepository interface {
	ResolveByID(id uuid.UUID) (someEntity SomeEntity, err error)
}

type SomeRepositoryMySQL struct {
	DB *infras.MysqlConn `inject:"db"`
}

func (r *SomeRepositoryMySQL) ResolveByID(id uuid.UUID) (someEntity SomeEntity, err error) {
	err = r.DB.Read.Get(&someEntity, fmt.Sprintf(querySample, id.String()))
	return
}
