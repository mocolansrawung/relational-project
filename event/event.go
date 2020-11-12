package event

type EventConsumer interface {
	Start()
}

type Event struct {
	Foo EventConsumer
}

func ProvideEvent(foo EventConsumer) Event {
	return Event{
		Foo: foo,
	}
}

// Start start all domains event consumer
func (e *Event) Start() {
	e.Foo.Start()
}
