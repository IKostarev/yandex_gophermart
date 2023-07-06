package config

import (
	"flag"
	"os"
	"yandex_gophermart/internal/models"
)

const (
	defaultAddr string = "localhost:8080"
)

func NewConfig() *models.Config {
	return Init(
		Addr(),
		Database(),
		Accrual(),
	)
}

func Database() models.Option {
	return func(p *models.Config) {
		flag.StringVar(&p.Database.ConnectionString, "d", "postgres://practicum:yandex@localhost:5432/postgres?sslmode=disable", "connection string for db")
		if envDBAddr := os.Getenv("DATABASE_URI"); envDBAddr != "" {
			p.Database.ConnectionString = envDBAddr
		}
	}
}

func Addr() models.Option {
	return func(p *models.Config) {
		flag.StringVar(&p.Server.Address, "a", defaultAddr, "address and port to run server")
		if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
			p.Server.Address = envRunAddr
		}
	}
}

func Accrual() models.Option {
	return func(p *models.Config) {
		flag.StringVar(&p.AccrualSystem.Address, "r", "", "address and port to run server")
		if envAccrualAddr := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualAddr != "" {
			p.AccrualSystem.Address = envAccrualAddr
		}
	}
}

func Init(opts ...models.Option) *models.Config {
	p := &models.Config{}

	for _, opt := range opts {
		opt(p)
	}

	flag.Parse()

	return p
}
