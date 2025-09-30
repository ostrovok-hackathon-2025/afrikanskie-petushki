package testhelper

import (
	"database/sql"
	"fmt"
	"path"
	"runtime"
	"testing"

	"github.com/jmoiron/sqlx"
	config2 "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/config"
	"github.com/pressly/goose"
)

const pgDriverName = "pgx"

func NewPostgre(t *testing.T) *sql.DB {
	t.Helper()
	config := config2.MustLoadConfig()
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.PostgresConfig.User,
		config.PostgresConfig.Password,
		config.PostgresConfig.Host,
		config.PostgresConfig.Port,
		config.PostgresConfig.Database,
	)
	postgresClient, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}
	return postgresClient
}

func NewPostgreSqlx(t *testing.T) *sqlx.DB {
	pgDB := NewPostgre(t)
	//_, filename, _, _ := runtime.Caller(0)
	//cfgPath := path.Dir(filename) + "/../../migrations/sql"
	//err := goose.Up(pgDB, cfgPath)
	//if err != nil {
	//	t.Fatal(err)
	//}

	return sqlx.NewDb(pgDB, pgDriverName)
}

func TearDownPostgreSqlx(t *testing.T, db *sqlx.DB) {
	_, filename, _, _ := runtime.Caller(0)
	cfgPath := path.Dir(filename) + "/../../migrations/sql"
	err := goose.DownTo(db.DB, cfgPath, 0)
	if err != nil {
		t.Fatal(err)
	}
}
