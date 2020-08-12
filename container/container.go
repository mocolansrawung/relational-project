package container

import (
	"github.com/facebookgo/inject"
)

type Service interface {
	OnStart()
	OnShutdown()
}

type ServiceRegistry struct {
	graph    inject.Graph
	objects  []*inject.Object
	services []interface{}
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

func (s *ServiceRegistry) Ready(svc interface{}) {
	switch obj := svc.(type) {
	case Service:
		obj.OnStart()
		s.services = append(s.services, svc)
	}
}

func (s *ServiceRegistry) Register(app string, svc interface{}) {
	s.Ready(svc)
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

func (s *ServiceRegistry) Shutdown() {
	if len(s.services) == 0 {
		return
	}

	for _, svc := range s.services {
		switch obj := svc.(type) {
		case Service:
			obj.OnShutdown()
		}
	}
}

func NewContainer() *ServiceRegistry {
	return &ServiceRegistry{}
}
