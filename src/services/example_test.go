package services

import (
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/evermos/boilerplate-go/mocks"
	"github.com/evermos/boilerplate-go/src/dto"
	"github.com/evermos/boilerplate-go/src/repositories"
)

func TestExampleService_Get(t *testing.T) {
	type fields struct {
		ExampleRepository repositories.ExampleContract
	}
	tests := []struct {
		name          string
		configureMock func(*mocks.MockExampleRepositoriesContract)
		want          dto.Example
		wantErr       bool
	}{
		{
			name: "status will return stable",
			configureMock: func(contract *mocks.MockExampleRepositoriesContract) {
				contract.EXPECT().Get().Return("stable", nil)
			},
			want:    dto.Example{Status: "stable"},
			wantErr: false,
		},
		{
			name: "get status will return error",
			configureMock: func(contract *mocks.MockExampleRepositoriesContract) {
				contract.EXPECT().Get().Return("", errors.New("error while get status"))
			},
			want:    dto.Example{Status: ""},
			wantErr: true,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repositoriesContract := mocks.NewMockExampleRepositoriesContract(ctrl)
			s := &ExampleService{
				ExampleRepository: repositoriesContract,
			}
			tt.configureMock(repositoriesContract)
			got, err := s.Get()
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}
