package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GRPC      GRPCServer
	Shopcart  Shopcart
	Warehouse Warehouse
	Postgres  Postgres
	Logger    Logger
}

type GRPCServer struct {
	Host string `env:"GRPC_HOST" env-required:"true"`
	Port string `env:"GRPC_PORT" env-required:"true"`
}

type Shopcart struct {
	Host string `env:"SHOPCART_HOST" env-required:"true"`
	Port string `env:"SHOPCART_PORT" env-required:"true"`
}

type Warehouse struct {
	Host string `env:"WAREHOUSE_HOST" env-required:"true"`
	Port string `env:"WAREHOUSE_PORT" env-required:"true"`
}

type Postgres struct {
	Host     string `env:"POSTGRES_HOST"     env-required:"true"`
	Port     string `env:"POSTGRES_PORT"     env-required:"true"`
	User     string `env:"POSTGRES_USER"     env-required:"true"`
	Name     string `env:"POSTGRES_DB"       env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	SSLMode  string `env:"POSTGRES_SSL_MODE" env-required:"true"`
}

type Logger struct {
	Level string `env:"LOGGER_LEVEL" env-default:"info"`
}

func MustLoad() *Config {
	cfg := &Config{}

	if err := cleanenv.ReadConfig(".env", cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return cfg
}
