package client

import (
	"github.com/Glawary/crypt/pkg/postgres"
)

func InitDB(cfg *postgres.Config) (*postgres.Instance, error) {
	var err error
	instance, err = postgres.NewPostgres(cfg)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

var instance *postgres.Instance

func GetDBInstance() *postgres.Instance {
	return instance
}
