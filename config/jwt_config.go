package config

type JWTConfig struct {
	SecretKey     string `env:"JWT_SECRET" yaml:"secret"`
	Issuer        string `env:"JWT_ISSUER" yaml:"issuer"`
	SigningMethod string `env:"JWT_SIGNING_METHOD" yaml:"signing_method"`
}
