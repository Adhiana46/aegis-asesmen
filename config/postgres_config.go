package config

type PostgresConfig struct {
	Host   string `env:"POSTGRES_HOST" yaml:"host"`
	Port   string `env:"POSTGRES_PORT" yaml:"port"`
	User   string `env:"POSTGRES_USER" yaml:"user"`
	Pass   string `env:"POSTGRES_PASS" yaml:"pass"`
	DbName string `env:"POSTGRES_DBNAME" yaml:"dbname"`
}
