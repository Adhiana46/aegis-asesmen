package publisher

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/Adhiana46/aegis-asesmen/pkg/kafka"
	"github.com/IBM/sarama"
)

var (
	ErrNoBrokers = errors.New("no Kafka bootstrap brokers defined")
	ErrNoConfig  = errors.New("no config specified")
)

type Publisher struct {
	producer sarama.SyncProducer
	config   *sarama.Config
	brokers  []string
}

func New(options ...Option) (*Publisher, error) {
	p := Publisher{
		brokers: []string{},
	}

	// apply options
	for _, opt := range options {
		opt(&p)
	}

	// if producer already set, imediately return
	if p.producer != nil {
		return &p, nil
	}

	// Validate configuration
	if len(p.brokers) == 0 {
		return nil, ErrNoBrokers
	}
	if p.config == nil {
		return nil, ErrNoConfig
	}

	// Create a new synchronous producer
	producer, err := sarama.NewSyncProducer(p.brokers, p.config)
	if err != nil {
		return nil, err
	}

	p.producer = producer

	return &p, nil
}

func (p *Publisher) Close() error {
	return p.producer.Close()
}

func (p *Publisher) Publish(e kafka.Event) error {
	topic := e.Topic()
	key := e.Key()
	payload, err := e.Payload()
	if err != nil {
		return err
	}

	msg := sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(payload),
	}

	partition, offset, err := p.producer.SendMessage(&msg)
	if err != nil {
		return err
	}
	slog.Info(fmt.Sprintf("Event published to topic: '%s' partition %d with offset %d", topic, partition, offset))

	return nil
}
