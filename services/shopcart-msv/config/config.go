package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	configPath = "./.env"
)

type Config struct {
	Storage    PostgresConfig
	GRPCServer GRPCServerConfig
	Logger     LoggerConfig
}

type PostgresConfig struct {
	User     string `env:"POSTGRES_USER"     env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	DB       string `env:"POSTGRES_DB"       env-required:"true"`
	Host     string `env:"POSTGRES_HOST"     env-required:"true"`
	Port     string `env:"POSTGRES_PORT"     env-required:"true"`
	SSLMode  string `env:"POSTGRES_SSL_MODE" env-required:"true"`
}

type GRPCServerConfig struct {
	Host string `env:"GRPC_HOST" env-required:"true"`
	Port string `env:"GRPC_PORT" env-required:"true"`
}

type LoggerConfig struct {
	Level string `default:"info" env:"LOG_LEVEL"`
}

func MustLoad() *Config {
	config := &Config{}

	err := cleanenv.ReadConfig(configPath, config)
	if err != nil {
		log.Fatalf("Error while loading config: %s", err)
	}

	return config
}
