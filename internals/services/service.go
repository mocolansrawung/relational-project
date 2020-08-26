package services

import "github.com/evermos/boilerplate-go/internals/dto"

type ExampleContract interface {
	Get() (dto.Example, error)
}
