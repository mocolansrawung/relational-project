package services

import "github.com/evermos/boilerplate-go/internal/dto"

type ExampleContract interface {
	Get() (dto.Example, error)
}
