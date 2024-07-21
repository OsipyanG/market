package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTP      HTTPServer
	Shopcart  Shopcart
	Order     Order
	Warehouse Warehouse
	Auth      Auth
	Logger    Logger
}

type HTTPServer struct {
	Host         string        `env:"HTTP_HOST"          env-required:"true"`
	Port         string        `env:"HTTP_PORT"          env-required:"true"`
	ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT"  env-default:"5"`
	WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT" env-default:"5"`
	IdleTimeout  time.Duration `env:"HTTP_IDLE_TIMEOUT"  env-default:"5"`
}

type Shopcart struct {
	Host string `env:"SHOPCART_HOST" env-required:"true"`
	Port string `env:"SHOPCART_PORT" env-required:"true"`
}

type Auth struct {
	Host string `env:"AUTH_HOST" env-required:"true"`
	Port string `env:"AUTH_PORT" env-required:"true"`
}

type Order struct {
	Host string `env:"ORDER_HOST" env-required:"true"`
	Port string `env:"ORDER_PORT" env-required:"true"`
}

type Warehouse struct {
	Host string `env:"WAREHOUSE_HOST" env-required:"true"`
	Port string `env:"WAREHOUSE_PORT" env-required:"true"`
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
