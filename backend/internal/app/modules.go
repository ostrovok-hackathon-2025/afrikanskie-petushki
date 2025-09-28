package app

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/config"
)

func initPostgresClient(config *config.PostgresConfig) *sqlx.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.User, config.Password, config.Host, config.Port, config.Database)
	postgresClient, err := sqlx.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}
	return postgresClient
}
