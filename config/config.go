package config

// Config grpc server configuration.
type Config struct {
	ServerAddress string

	ServerCA          string
	ServerCertificate string
	ServerKey         string

	ClientCA          string
	ClientCertificate string
	ClientKey         string
}

func NewConfig() Config {
	return Config{}
}

