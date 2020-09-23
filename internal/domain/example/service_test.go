package example_test

import (
	"testing"
	"time"

	"github.com/evermos/boilerplate-go/internal/domain/example"
	example_mock "github.com/evermos/boilerplate-go/internal/domain/example/mock"
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

func TestExampleService(t *testing.T) {

	t.Run("resolveByID", func(t *testing.T) {
		tests := []struct {
			name      string
			entityID  uuid.UUID
			setupMock func(*example_mock.MockSomeRepository, uuid.UUID, example.SomeEntity, error)
			returns   *example.SomeEntity
			err       error
		}{
			{
				name:     "default",
				entityID: getRandomUUID(),
				setupMock: func(mockRepo *example_mock.MockSomeRepository, id uuid.UUID, ent example.SomeEntity, err error) {
					mockRepo.EXPECT().ResolveByID(id).Return(ent, err)
				},
				returns: &example.SomeEntity{
					Name:      "John Doe",
					Status:    example.SomeEntityStatusActive,
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
				mockRepo := example_mock.NewMockSomeRepository(ctrl)
				s := &example.SomeServiceImpl{
					SomeRepository: mockRepo,
				}
				test.setupMock(mockRepo, test.entityID, *test.returns, test.err)
				got, err := s.ResolveByID(test.entityID)

				assert.Equal(t, test.err, err)
				assert.Equal(t, test.returns.Name, got.Name)
			})
		}
	})
}
