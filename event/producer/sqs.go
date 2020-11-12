package producer

import (
	"encoding/json"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/event/model"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
)

type SQSConfig struct {
	Config configs.Config
}

func (m *SQSConfig) Retrieve() (credentials.Value, error) {
	return credentials.Value{
		AccessKeyID:     m.Config.Event.Producer.SQS.AccessKeyID,
		SecretAccessKey: m.Config.Event.Producer.SQS.SecretAccessKey,
	}, nil
}

func (m *SQSConfig) IsExpired() bool {
	return false
}

func createSqsConfig(config *configs.Config) (*session.Session, error) {
	sqsConfig := SQSConfig{Config: *config}
	return session.NewSession(&aws.Config{
		Region:      &config.Event.Producer.SQS.Region,
		Credentials: credentials.NewCredentials(&sqsConfig),
		MaxRetries:  aws.Int(config.Event.Producer.SQS.MaxRetries),
	})
}

type SQSType string

const (
	Standard SQSType = "standard"
	Fifo     SQSType = "fifo"
)

type SQSProducer struct {
	config *configs.Config
	sqs    *sqs.SQS
}

// NewProducer create object from Producer
func NewSqsProducer(config *configs.Config) *SQSProducer {
	sess, err := createSqsConfig(config)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed creating sqs config")
	}
	return &SQSProducer{config: config, sqs: sqs.New(sess)}
}

// Send is function to produce event to Queue
func (p *SQSProducer) Send(event model.EventWrapper, url string) {
	value, err := json.Marshal(event)
	if err != nil {
		log.Err(err).Msgf("Error converting event to json: %v", err)
		return
	}

	message := &sqs.SendMessageInput{
		MessageBody:    aws.String(string(value)),
		QueueUrl:       aws.String(url),
		DelaySeconds:   aws.Int64(p.config.Event.Producer.SQS.DelayPeriodSeconds),
		MessageGroupId: aws.String(url),
	}

	fifo := strings.Split(url, ".")
	if fifo[len(fifo)-1] == string(Fifo) {
		id, _ := uuid.NewV4()
		message.MessageDeduplicationId = aws.String(id.String())
	}

	go p.send(message)
}

func (p *SQSProducer) send(msg *sqs.SendMessageInput) {
	resp, err := p.sqs.SendMessage(msg)
	switch {
	case err != nil:
		log.Err(err).Msgf("Error sending message : %v %s", err, *msg.MessageBody)
	case err == nil:
		log.Info().Msgf("Success sending message : %v", resp.String())
	}
}
