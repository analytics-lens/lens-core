package config

type ServerConfig struct {
	Port    string
	Passkey string
}

type DatabaseConfig struct {
	SSL      string
	Name     string
	Host     string
	Port     string
	Password string
	Username string
}

type OpenAIConfig struct {
	SecretKey string
}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	OpenAI   OpenAIConfig
}
