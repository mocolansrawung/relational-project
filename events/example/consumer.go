package example

import (
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/events/consumer"
	"github.com/evermos/boilerplate-go/internal/domain/example"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"

	"github.com/cenkalti/backoff"
	"github.com/nsqio/go-nsq"
)

type EventConsumer struct {
	nsq         consumer.NsqEventConsumer
	Config      *configs.Config     `inject:"config"`
	SomeService example.SomeService `inject:"example.someService"`
}

func (e *EventConsumer) OnStart() {
	e.nsq = consumer.NsqEventConsumer{
		Config:        e.Config,
		Topic:         "evm.example-service.status-created",
		Channel:       "evm.example-service.payment-webhook",
		HandleMessage: e.processEvent,
	}

	if e.Config.EnableExampleConsumer {
		log.Info().Msg("Starting example consumer...")
		go e.nsq.Start()
	}
}

func (e *EventConsumer) processEvent(message *nsq.Message) error {
	log.Info().Msg("Start processing event...")

	backoffWithMaxRetry := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), e.Config.ConsumerBackoffMaxRetry)
	process := func() error {
		randomID, err := uuid.NewV4()
		_, err = e.SomeService.ResolveByID(randomID)
		if err != nil {
			log.Err(err).Msgf("Failed to resolve id %v", randomID)
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
	log.Info().Msg("Shutting down example consumer...")
	e.nsq.Stop()
}
