package foobarbaz

import (
	"encoding/json"

	"github.com/cenkalti/backoff/v4"
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/event/consumer"
	"github.com/evermos/boilerplate-go/event/model"
	"github.com/evermos/boilerplate-go/internal/domain/foobarbaz"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
)

type FooBarConsumerImpl struct {
	Config     *configs.Config
	FooService foobarbaz.FooService
	Consumer   consumer.Consumer
}

// ProvideFooBarConsumerImpl is the provider for this consumer.
func ProvideFooBarConsumerImpl(config *configs.Config, fooService foobarbaz.FooService) *FooBarConsumerImpl {
	f := FooBarConsumerImpl{}
	f.Config = config
	f.FooService = fooService

	sqsConsumer := consumer.NewSQSConsumer(config)
	sqsConsumer.Process = f.processEvent
	f.Consumer = sqsConsumer

	return &f
}

// Start running listener
func (c *FooBarConsumerImpl) Start() {
	go c.Consumer.Listen(c.Config.Event.Consumer.SQS.TopicURLs.FooBar)
}

func (c *FooBarConsumerImpl) processEvent(value []byte) error {
	evt := model.EventWrapper{}
	err := json.Unmarshal(value, &evt)
	if err != nil {
		return err
	}

	foo := foobarbaz.FooRequestFormat{}
	err = json.Unmarshal(evt.Data.Value, &foo)
	if err != nil {
		return err
	}

	err = c.updateProcessFoo(foo)
	if err != nil {
		return err
	}

	return nil
}

func (c *FooBarConsumerImpl) updateProcessFoo(foo foobarbaz.FooRequestFormat) error {
	backoffWithMaxRetry := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), c.Config.Event.Consumer.SQS.MaxRetriesConsume)

	// example process
	// TODO: use proper example function
	ID, _ := uuid.NewV4()
	userID, _ := uuid.NewV4()
	consume := func() error {
		foo, err := c.FooService.Update(ID, foo, userID)
		if err != nil {
			return err
		}

		log.Info().Msgf("Foo updated: %#v", foo)
		return nil
	}

	err := backoff.Retry(consume, backoffWithMaxRetry)
	if err != nil {
		log.Err(err).Msgf("err consuming process :")
		return err
	}
	return nil
}
