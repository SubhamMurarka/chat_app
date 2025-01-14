package db

import (
	"context"
	"fmt"
	"log"

	"github.com/SubhamMurarka/chat_app/User/Config"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type SQLDatabase struct {
	db *pgxpool.Pool
}

func getConfig(url string) *pgxpool.Pool {
	pgxConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatalf("Error getting config object : %v", err)
		return nil
	}

	stdlib.RegisterConnConfig(pgxConfig.ConnConfig)

	pgxConfig.MinConns = 5
	pgxConfig.MaxConns = 20

	pgxpool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		log.Fatalf("Error getting config object : %v", err)
		return nil
	}

	return pgxpool
}

func NewSQLDatabase() (*SQLDatabase, error) {
	commonstring := fmt.Sprintf("postgres://%s:%s@", Config.Config.PostgresUser, Config.Config.PostgresPassword)

	pgPool := GetConfig(commonstring + Config.Config.PostgresHost3 + ":" + Config.Config.PostgresPort + "/user?sslmode=disable")

	return &SQLDatabase{db: pgPool}, nil
}

func (d *SQLDatabase) GetDB() *pgxpool.Pool {
	return d.db
}

func (d *SQLDatabase) Close() {
	d.db.Close()
}
