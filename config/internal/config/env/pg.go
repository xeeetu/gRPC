package env

import (
	"errors"
	"os"

	"github.com/xeeetu/gRPC/config/internal/config"
)

const (
	dsnEnvName = "PG_DSN"
)

var _ config.PGConfig = (*pgConfig)(nil)

type pgConfig struct {
	dsn string
}

func NewPGConfig() (*pgConfig, error) {
	dsn := os.Getenv(dsnEnvName)
	if dsn == "" {
		return nil, errors.New("pg dsn not found")
	}
	return &pgConfig{dsn: dsn}, nil
}

func (cfg *pgConfig) DSN() string {
	return cfg.dsn
}
