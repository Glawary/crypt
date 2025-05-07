package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Instance struct {
	db  *sqlx.DB
	cfg *Config
}

func NewPostgres(cfg *Config) (*Instance, error) {
	instance := &Instance{
		cfg: cfg,
	}

	conn, err := sqlx.Open(cfg.Driver, cfg.FormDSN())
	if err != nil {
		return nil, fmt.Errorf("could not open database connection: %v", err)
	}

	err = conn.PingContext(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not establish database connection: %v", err)
	}
	instance.db = conn

	return instance, nil
}

func (rec *Instance) GetSqlxDB() *sqlx.DB {
	return rec.db
}
