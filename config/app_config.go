package config

type AppConfig struct {
	Name    string `env-required:"true" env:"APP_NAME" yaml:"name"`
	Version string `env-required:"true" env:"APP_VERSION" yaml:"version"`
}
