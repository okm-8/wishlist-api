package config

type HttpServerConfig struct {
	address string
}

func NewHttpServerConfig(address string) *HttpServerConfig {
	return &HttpServerConfig{
		address: address,
	}
}

func (config *HttpServerConfig) Address() string {
	return config.address
}

type HttpConfig struct {
	publicServer  *HttpServerConfig
	privateServer *HttpServerConfig
}

func NewHttpConfig(
	publicServer *HttpServerConfig,
	privateServer *HttpServerConfig,
) *HttpConfig {
	return &HttpConfig{
		publicServer:  publicServer,
		privateServer: privateServer,
	}
}

func (config *HttpConfig) PublicServer() *HttpServerConfig {
	return config.publicServer
}

func (config *HttpConfig) PrivateServer() *HttpServerConfig {
	return config.privateServer
}
