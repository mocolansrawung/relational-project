package example

import (
	"time"

	"github.com/evermos/boilerplate-go/shared/nuuid"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
)

// SomeEntity is a sample entity model
type SomeEntity struct {
	ID        uuid.UUID   `db:"id" json:"id"`
	Name      string      `db:"name" json:"name"`
	Created   time.Time   `db:"created" json:"created"`
	CreatedBy uuid.UUID   `db:"created_by" json:"createdBy"`
	Updated   null.Time   `db:"updated" json:"updated,omitempty"`
	UpdatedBy nuuid.NUUID `db:"updated_by" json:"updatedBy,omitempty"`
}
