package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"

	"github.com/Glawary/crypt/pkg/grpc"
	"github.com/Glawary/crypt/pkg/http"
	"github.com/Glawary/crypt/pkg/postgres"
)

type (
	Config struct {
		GRPCServer *grpc.GRPCConfig
		HttpServer *http.HttpConfig
		DB         *postgres.Config
	}
)

func New(path string) (*Config, error) {
	cfg := &Config{
		GRPCServer: &grpc.GRPCConfig{},
		HttpServer: &http.HttpConfig{},
		DB:         &postgres.Config{},
	}

	err := godotenv.Overload(path)
	if err != nil {
		return nil, err
	}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
