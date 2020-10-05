package foobarbaz

import (
	"encoding/json"
	"time"

	"github.com/evermos/boilerplate-go/shared/nuuid"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
)

// FooStatus indicates the status of Foo
type FooStatus string

const (
	// FooStatusActive indicates an active Foo
	FooStatusActive FooStatus = "active"
	// FooStatusInactive indicates an inactive Foo
	FooStatusInactive FooStatus = "inactive"
)

// Foo is a sample entity model
type Foo struct {
	ID        uuid.UUID   `db:"id" json:"id"`
	Name      string      `db:"name" json:"name"`
	Status    FooStatus   `db:"status" json:"status"`
	Created   time.Time   `db:"created" json:"created"`
	CreatedBy uuid.UUID   `db:"created_by" json:"createdBy"`
	Updated   null.Time   `db:"updated" json:"updated,omitempty"`
	UpdatedBy nuuid.NUUID `db:"updated_by" json:"updatedBy,omitempty"`
}

// FooFormat represents a Foo's standard formatting for JSON serializing.
type FooFormat struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Status    FooStatus  `json:"status"`
	Created   time.Time  `json:"created"`
	CreatedBy uuid.UUID  `json:"createdBy"`
	Updated   *time.Time `json:"updated,omitempty"`
	UpdatedBy *uuid.UUID `json:"updatedBy,omitempty"`
}

// MarshalJSON overrides the standard JSON formatting.
func (f *Foo) MarshalJSON() ([]byte, error) {
	format := FooFormat{
		ID:        f.ID,
		Name:      f.Name,
		Status:    f.Status,
		Created:   f.Created,
		CreatedBy: f.CreatedBy,
		Updated:   f.Updated.Ptr(),
		UpdatedBy: f.UpdatedBy.Ptr(),
	}

	return json.Marshal(format)
}
