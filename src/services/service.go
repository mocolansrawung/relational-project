package services

import (
	"fmt"

	"github.com/evermos/boilerplate-go/src/repositories"

	"github.com/evermos/boilerplate-go/configs"
)

type Service struct {
	Config     *configs.Config          `inject:"config"`
	Repository *repositories.Repository `inject:"repo"`
}

func (s *Service) TestService() {
	s.Repository.Get()
	fmt.Println("Get From Service")
}
