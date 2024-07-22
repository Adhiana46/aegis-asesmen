package config

type KafkaConfig struct {
	ClientId      string   `env:"KAFKA_CLIENT_ID" yaml:"client_id"`
	Username      string   `env:"KAFKA_USERNAME" yaml:"username"`
	Password      string   `env:"KAFKA_PASSWORD" yaml:"password"`
	Brokers       []string `env-separator:"," env:"KAFKA_BROKERS" yaml:"brokers"`
	SASLMechanism string   `env:"KAFKA_SASL_MECHANISM" yaml:"sasl_mechanism"`
}
