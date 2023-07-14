package config

import (
	"flag"
	"os"
)

const (
	defaultAddr string = "localhost:8080"
)

type Option func(params *Config)

type Config struct {
	Server struct {
		Address string
	}
	Database struct {
		ConnectionString string
	}
	AccrualSystem struct {
		Address string
	}
}

func NewConfig() *Config {
	return Init(
		addr(),
		database(),
		accrual(),
	)
}

func database() Option {
	return func(p *Config) {
		flag.StringVar(&p.Database.ConnectionString, "d", "postgres://practicum:yandex@localhost:5432/postgres?sslmode=disable", "connection string for db")
		if envDBAddr := os.Getenv("DATABASE_URI"); envDBAddr != "" {
			p.Database.ConnectionString = envDBAddr
		}
	}
}

func addr() Option {
	return func(p *Config) {
		flag.StringVar(&p.Server.Address, "a", defaultAddr, "address and port to run server")
		if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
			p.Server.Address = envRunAddr
		}
	}
}

func accrual() Option {
	return func(p *Config) {
		flag.StringVar(&p.AccrualSystem.Address, "r", "", "address and port to run server")
		if envAccrualAddr := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualAddr != "" {
			p.AccrualSystem.Address = envAccrualAddr
		}
	}
}

func Init(opts ...Option) *Config {
	p := &Config{}

	for _, opt := range opts {
		opt(p)
	}

	flag.Parse()

	return p
}
