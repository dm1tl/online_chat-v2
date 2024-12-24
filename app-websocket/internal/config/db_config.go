package config

import (
	"errors"
	"os"
)

type DBConfig interface {
	DSN() string
}

const (
	pg_dsn = "PG_DSN"
)

type envDBConfig struct {
	dsn string
}

func NewDBConfig() (*envDBConfig, error) {
	dsn := os.Getenv(pg_dsn)
	if len(dsn) == 0 {
		return nil, errors.New("pg_dsn is empty")
	}
	return &envDBConfig{
		dsn: dsn,
	}, nil
}

func (p *envDBConfig) DSN() string {
	return p.dsn
}
