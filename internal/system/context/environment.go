package context

import "api/internal/service/integration/environment"

type Environment struct {
	AppSecret string `env:"APP_SECRET" required:"true"`

	PostgresDsn string `env:"POSTGRES_DSN" required:"true"`
	RedisDsn    string `env:"REDIS_DSN" required:"true"`

	MigrationsPath string `env:"MIGRATIONS_PATH" default:"migrations"`

	PublicAddress  string `env:"PUBLIC_ADDRESS" default:"localhost:8080"`
	PrivateAddress string `env:"PRIVATE_ADDRESS" default:"localhost:8081"`

	LogLevel   string `env:"LOG_LEVEL" default:"INFO"`
	LogOutput  string `env:"LOG_OUTPUT" default:"stdout"`
	ColoredLog bool   `env:"COLORED_LOG" default:"false"`
}

func NewEnvironment() (*Environment, error) {
	env := &Environment{}

	return env, environment.Read(env, ".env", ".env.local")
}
