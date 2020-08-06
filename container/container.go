package container

import (
	"github.com/facebookgo/inject"
)

type ServiceRegistry struct {
	graph   inject.Graph
	objects []*inject.Object
}

func (s *ServiceRegistry) bind() error {
	for _, obj := range s.objects {
		err := s.graph.Provide(obj)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *ServiceRegistry) Register(app string, svc interface{}) {
	s.objects = append(s.objects, &inject.Object{Value: svc, Name: app})
}

func (s *ServiceRegistry) Start() error {
	err := s.bind()
	if err != nil {
		return err
	}

	err = s.graph.Populate()
	if err != nil {
		return err
	}

	return nil
}

func NewContainer() *ServiceRegistry {
	return &ServiceRegistry{}
}
