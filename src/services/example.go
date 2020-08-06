package services

import (
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/src/dto"
	"github.com/evermos/boilerplate-go/src/repositories"
)

type ExampleService struct {
	Config            *configs.Config                 `inject:"config"`
	ExampleRepository *repositories.ExampleRepository `inject:"repository.example"`
}

func (s *ExampleService) Get() dto.Example {
	return dto.Example{Status: s.ExampleRepository.Get()}
}
