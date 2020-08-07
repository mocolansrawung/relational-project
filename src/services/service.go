package services

import "github.com/evermos/boilerplate-go/src/dto"

type ExampleContract interface {
	Get() (dto.Example, error)
}
