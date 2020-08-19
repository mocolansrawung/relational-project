package producer

import (
	"fmt"
	"log"
	"time"

	"github.com/evermos/boilerplate-go/configs"

	"github.com/nsqio/go-nsq"
)

type Producer struct {
	Config   *configs.Config `inject:"config"`
	Producer *nsq.Producer
}

func createNsqConfig(cfg *configs.Config) *nsq.Config {
	config := nsq.NewConfig()
	config.DialTimeout = time.Duration(cfg.ProducerDialTimeout) * time.Second
	config.BackoffMultiplier = time.Duration(cfg.ProducerRetryBackoff) * time.Second
	config.BackoffStrategy.Calculate(cfg.ProducerBackoffMaxRetry)

	return config
}

func (p *Producer) OnStart() {
	log.Println("Initiliazing Producer...")
	nsqConfig := createNsqConfig(p.Config)
	producer, err := nsq.NewProducer(fmt.Sprintf("%s:%s", p.Config.NsqHost, p.Config.NsqPort), nsqConfig)
	if err != nil {
		log.Fatalln(err)
	}
	p.Producer = producer
}

func (p *Producer) Emit(message, topic string) error {
	err := p.Producer.Publish(topic, []byte(message))
	if err != nil {
		return err
	}
	return nil
}

func (p *Producer) OnShutdown() {
	log.Println("Stopping Producer...")
	p.Producer.Stop()
}
