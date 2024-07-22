package publisher

import "github.com/IBM/sarama"

type Option func(*Publisher)

func WithBrokers(brokers []string) Option {
	return func(p *Publisher) {
		p.brokers = brokers
	}
}

func WithConfig(config *sarama.Config) Option {
	return func(p *Publisher) {
		p.config = config
	}
}

func WithProducer(producer sarama.SyncProducer) Option {
	return func(p *Publisher) {
		p.producer = producer
	}
}
