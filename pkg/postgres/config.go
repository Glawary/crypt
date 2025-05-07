package postgres

import (
	"fmt"
)

type Config struct {
	Port    uint   `env:"DB_PORT,required"`
	Driver  string `env:"DB_DRIVER" envDefault:"postgres"`
	Host    string `env:"DB_HOST,required"`
	Pass    string `env:"DB_PASS,required"`
	User    string `env:"DB_USER,required"`
	Name    string `env:"DB_NAME,required"`
	SSLMode string `env:"DB_SSL_MODE,required"`
}

func (rec *Config) FormDSN() string {
	return fmt.Sprintf(
		"%s://%s:%s@%s:%d/%s?sslmode=%s",
		rec.Driver,
		rec.User,
		rec.Pass,
		rec.Host,
		rec.Port,
		rec.Name,
		rec.SSLMode,
	)
}
