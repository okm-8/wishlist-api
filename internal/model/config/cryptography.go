package config

type CryptographyConfig struct {
	secret string
}

func NewCryptographyConfig(secret string) *CryptographyConfig {
	return &CryptographyConfig{
		secret: secret,
	}
}

func (env *CryptographyConfig) Secret() string {
	return env.secret
}
