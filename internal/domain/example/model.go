package example

import (
	"encoding/json"
	"time"

	"github.com/evermos/boilerplate-go/shared/nuuid"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
)

// SomeEntityStatus indicates the status of SomeEntity
type SomeEntityStatus string

const (
	// SomeEntityStatusActive indicates an active SomeEntity
	SomeEntityStatusActive SomeEntityStatus = "active"
	// SomeEntityStatusInactive indicates an inactive SomeEntity
	SomeEntityStatusInactive SomeEntityStatus = "inactive"
)

// SomeEntity is a sample entity model
type SomeEntity struct {
	ID        uuid.UUID        `db:"id" json:"id"`
	Name      string           `db:"name" json:"name"`
	Status    SomeEntityStatus `db:"status" json:"status"`
	Created   time.Time        `db:"created" json:"created"`
	CreatedBy uuid.UUID        `db:"created_by" json:"createdBy"`
	Updated   null.Time        `db:"updated" json:"updated,omitempty"`
	UpdatedBy nuuid.NUUID      `db:"updated_by" json:"updatedBy,omitempty"`
}

// SomeEntityFormat represents a SomeEntity's standard formatting for JSON serializing.
type SomeEntityFormat struct {
	ID        uuid.UUID        `json:"id"`
	Name      string           `json:"name"`
	Status    SomeEntityStatus `json:"status"`
	Created   time.Time        `json:"created"`
	CreatedBy uuid.UUID        `json:"createdBy"`
	Updated   *time.Time       `json:"updated,omitempty"`
	UpdatedBy *uuid.UUID       `json:"updatedBy,omitempty"`
}

// MarshalJSON overrides the standard JSON formatting.
func (se *SomeEntity) MarshalJSON() ([]byte, error) {
	format := SomeEntityFormat{
		ID:        se.ID,
		Name:      se.Name,
		Status:    se.Status,
		Created:   se.Created,
		CreatedBy: se.CreatedBy,
		Updated:   se.Updated.Ptr(),
		UpdatedBy: se.UpdatedBy.Ptr(),
	}

	return json.Marshal(format)
}
