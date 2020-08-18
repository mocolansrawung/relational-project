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

func (s *ServiceRegistry) Ready(svc interface{}) {
	switch obj := svc.(type) {
	case Service:
		obj.OnStart()
	}
}

func (s *ServiceRegistry) Register(app string, svc interface{}) {
	err := s.graph.Provide(&inject.Object{Value: svc, Name: app})
	if err != nil {
		panic(err.Error())
	}
	s.services = append(s.services, svc)
}

func (s *ServiceRegistry) Start() error {
	err := s.graph.Populate()
	if err != nil {
		return err
	}

	for _, svc := range s.services {
		s.Ready(svc)
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
