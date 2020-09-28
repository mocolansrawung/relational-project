package foobarbaz_test

import (
	"testing"
	"time"

	"github.com/evermos/boilerplate-go/internal/domain/foobarbaz"
	foobarbaz_mock "github.com/evermos/boilerplate-go/internal/domain/foobarbaz/mock"
	"github.com/evermos/boilerplate-go/shared/nuuid"
	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
)

func getRandomUUID() uuid.UUID {
	id, _ := uuid.NewV4()
	return id
}

func TestFooService(t *testing.T) {

	t.Run("resolveByID", func(t *testing.T) {
		tests := []struct {
			name      string
			entityID  uuid.UUID
			setupMock func(*foobarbaz_mock.MockFooRepository, uuid.UUID, foobarbaz.Foo, error)
			returns   *foobarbaz.Foo
			err       error
		}{
			{
				name:     "default",
				entityID: getRandomUUID(),
				setupMock: func(mockRepo *foobarbaz_mock.MockFooRepository, id uuid.UUID, ent foobarbaz.Foo, err error) {
					mockRepo.EXPECT().ResolveByID(id).Return(ent, err)
				},
				returns: &foobarbaz.Foo{
					Name:      "John Doe",
					Status:    foobarbaz.FooStatusActive,
					Created:   time.Now(),
					CreatedBy: getRandomUUID(),
					Updated:   null.TimeFrom(time.Now()),
					UpdatedBy: nuuid.From(getRandomUUID()),
				},
				err: nil,
			},
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				mockRepo := foobarbaz_mock.NewMockFooRepository(ctrl)
				s := &foobarbaz.FooServiceImpl{
					FooRepository: mockRepo,
				}
				test.setupMock(mockRepo, test.entityID, *test.returns, test.err)
				got, err := s.ResolveByID(test.entityID)

				assert.Equal(t, test.err, err)
				assert.Equal(t, test.returns.Name, got.Name)
			})
		}
	})
}
