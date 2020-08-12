package example

import (
	"github.com/nsqio/go-nsq"
	"log"

	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/events/consumer"
)

type EventConsumer struct {
	Config *configs.Config `inject:"config"`
	nsq    consumer.NsqEventConsumer
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
	return nil
}

func (e *EventConsumer) OnShutDown() {
	log.Println("shutting down nsq")
	e.nsq.Stop()
}
