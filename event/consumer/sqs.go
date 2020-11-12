package consumer

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/evermos/boilerplate-go/configs"
	"github.com/rs/zerolog/log"
)

type Process func(e []byte) error

type SQSConfig struct {
	Config configs.Config
}

func (m *SQSConfig) IsExpired() bool {
	return false
}

func (m *SQSConfig) Retrieve() (credentials.Value, error) {
	return credentials.Value{
		AccessKeyID:     m.Config.Event.Consumer.SQS.AccessKeyID,
		SecretAccessKey: m.Config.Event.Consumer.SQS.SecretAccessKey,
	}, nil
}

func createSQSConfig(config *configs.Config) (*session.Session, error) {
	sqsConfig := SQSConfig{Config: *config}
	return session.NewSession(&aws.Config{
		Region:      &config.Event.Consumer.SQS.Region,
		Credentials: credentials.NewCredentials(&sqsConfig),
		MaxRetries:  aws.Int(config.Event.Consumer.SQS.MaxRetries),
	})
}

type SQSConsumer struct {
	Process Process
	config  *configs.Config
	sqs     *sqs.SQS
}

// NewConsumer create object Consumer
func NewSQSConsumer(config *configs.Config) *SQSConsumer {
	sess, err := createSQSConfig(config)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed creating sqs config")
	}
	return &SQSConsumer{config: config, sqs: sqs.New(sess)}
}

// Listen is a function to listen new message from sqs queue
func (p *SQSConsumer) Listen(url string) {
	log.Info().Msgf("Start listen url : %v", url)

	for {
		receiveResp, err := p.sqs.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(url),
			MaxNumberOfMessages: aws.Int64(p.config.Event.Consumer.SQS.MaxMessage),
			WaitTimeSeconds:     aws.Int64(p.config.Event.Consumer.SQS.WaitTimeSeconds),
		})
		if err != nil {
			log.Err(err).Msgf("error receiving message :")
		}

		for _, message := range receiveResp.Messages {
			err := p.Process([]byte(*message.Body))
			if err != nil {
				log.Err(err).Msgf("error process message :")
			}

			err = p.deleteMessage(message, url)
			if err != nil {
				log.Err(err).Msgf("error delete message :")
			}
		}

		time.Sleep(time.Duration(p.config.Event.Consumer.SQS.IntervalPeriodSeconds) * time.Second)
	}
}

func (p *SQSConsumer) deleteMessage(msg *sqs.Message, url string) error {
	output, err := p.sqs.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &url,
		ReceiptHandle: msg.ReceiptHandle,
	})
	if err != nil {
		log.Err(err).Msgf("Error DeleteMessage: %v, Message output: %v", err, output)
		return err
	}
	return nil
}
