package example

import (
	"log"

	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/events/consumer"
	"github.com/evermos/boilerplate-go/src/services"

	"github.com/cenkalti/backoff"
	"github.com/nsqio/go-nsq"
)

type EventConsumer struct {
	nsq            consumer.NsqEventConsumer
	Config         *configs.Config          `inject:"config"`
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
		log.Println("Starting example consumer...")
		go e.nsq.Start()
	}
}

func (e *EventConsumer) processEvent(message *nsq.Message) error {
	log.Println("Start processing event...")

	backoffWithMaxRetry := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), e.Config.ConsumerBackoffMaxRetry)
	process := func() error {
		_, err := e.ExampleService.Get()
		if err != nil {
			log.Println("err : ", err.Error())
			return err
		}

		return nil
	}

	err := backoff.Retry(process, backoffWithMaxRetry)
	if err != nil {
		return err
	}

	return nil
}

func (e *EventConsumer) OnShutdown() {
	log.Println("Shutting down example consumer...")
	e.nsq.Stop()
}
