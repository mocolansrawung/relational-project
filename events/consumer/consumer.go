package consumer

import (
	"fmt"
	"sync"

	"github.com/evermos/boilerplate-go/configs"
	"github.com/nsqio/go-nsq"
	"github.com/rs/zerolog/log"
)

var (
	waitGroup sync.WaitGroup
)

type NsqEventConsumer struct {
	Config        *configs.Config
	Topic         string
	Channel       string
	HandleMessage func(message *nsq.Message) error
	nsq           *nsq.Config
	consumer      *nsq.Consumer
}

func (e *NsqEventConsumer) Start() {
	log.Info().Msgf("Starting nsq consumer : topic %v - channel %v", e.Topic, e.Channel)
	e.nsq = nsq.NewConfig()
	consumer, err := nsq.NewConsumer(e.Topic, e.Channel, e.nsq)
	if err != nil {
		log.Fatal().Err(err).Msg("failed creating consumer")
		return
	}
	e.consumer = consumer

	err = e.Consume()
	if err != nil {
		log.Error().Err(err).Msg("failed consuming event")
	}
}

func (e *NsqEventConsumer) Consume() error {
	e.consumer.AddHandler(nsq.HandlerFunc(e.HandleMessage))
	nsqdAddress := fmt.Sprintf("%s:%s", e.Config.NsqHost, e.Config.NsqPort)
	err := e.consumer.ConnectToNSQD(nsqdAddress)
	if err != nil {
		log.Err(err).Msgf("Failed to connect to Nsq %v", nsqdAddress)
		return err
	}
	waitGroup.Wait()
	return nil
}

func (e *NsqEventConsumer) Stop() {
	log.Info().Msg("stoping consumer...")
	e.consumer.Stop()
}
