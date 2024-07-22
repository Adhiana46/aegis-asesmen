package config

type HttpServerConfig struct {
	Host string `env-required:"true" env:"HTTP_SERVER_HOST" yaml:"host"`
	Port string `env-required:"true" env:"HTTP_SERVER_PORT" yaml:"port"`
}
