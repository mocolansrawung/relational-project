package example

import (
	"github.com/evermos/boilerplate-go/src/services"
	"log"

	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/events/consumer"

	"github.com/nsqio/go-nsq"
)

type EventConsumer struct {
	Config         *configs.Config `inject:"config"`
	nsq            consumer.NsqEventConsumer
	ExampleService services.ExampleContract `inject:"service.example"`
}

func (e *EventConsumer) OnStart() {
	e.nsq = consumer.NsqEventConsumer{
		Config:        e.Config,
		Topic:         "evm.example-service.status-created",
		Channel:       "evm.example-service.payment-webhook",
		HandleMessage: e.processEvent,
	}

	if e.Config.EnableExampleConsumer {
		go e.nsq.Start()
	}
}

func (e *EventConsumer) processEvent(message *nsq.Message) error {
	log.Println("Start processing event...")
	_, err := e.ExampleService.Get()
	if err != nil {
		log.Println("err : ", err.Error())
		return err
	}
	return nil
}

func (e *EventConsumer) OnShutDown() {
	log.Println("shutting down nsq")
	e.nsq.Stop()
}
