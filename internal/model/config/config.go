package config

type Config struct {
	cryptography *CryptographyConfig
	pgx          *PgxConfig
	migrations   *MigrationsConfig
	http         *HttpConfig
	logger       *LoggerConfig
}

func NewConfig(
	cryptography *CryptographyConfig,
	pgx *PgxConfig,
	migrations *MigrationsConfig,
	http *HttpConfig,
	logger *LoggerConfig,
) *Config {
	return &Config{
		cryptography: cryptography,
		pgx:          pgx,
		migrations:   migrations,
		http:         http,
		logger:       logger,
	}
}

func (config *Config) Cryptography() *CryptographyConfig {
	return config.cryptography
}

func (config *Config) Pgx() *PgxConfig {
	return config.pgx
}

func (config *Config) Migrations() *MigrationsConfig {
	return config.migrations
}

func (config *Config) Http() *HttpConfig {
	return config.http
}

func (config *Config) LoggerConfig() *LoggerConfig {
	return config.logger
}
