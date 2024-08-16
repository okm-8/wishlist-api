package config

type PgxConfig struct {
	postgresDsn string
}

func NewPgxConfig(postgresDsn string) *PgxConfig {
	return &PgxConfig{postgresDsn: postgresDsn}
}

func (env *PgxConfig) PostgresDsn() string {
	return env.postgresDsn
}
