package services

import (
	"log"

	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/src/dto"
	"github.com/evermos/boilerplate-go/src/repositories"
)

type ExampleService struct {
	Config            *configs.Config              `inject:"config"`
	ExampleRepository repositories.ExampleContract `inject:"repository.example"`
}

func (s *ExampleService) OnStart() {
	log.Println("Start initializing...")
}

func (s *ExampleService) OnShutdown() {
	log.Println("Shutdown...")
}

func (s *ExampleService) Get() (dto.Example, error) {
	status, err := s.ExampleRepository.Get()
	if err != nil {
		return dto.Example{}, err
	}
	return dto.Example{Status: status}, nil
}
