package db

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type SQLDatabase struct {
	db *sql.DB
}

func NewSQLDatabase() (*SQLDatabase, error) {
	db, err := sql.Open("postgres", "postgres://root:password@localhost:5435/go-chat?sslmode=disable")
	if err != nil {
		return nil, err
	}

	return &SQLDatabase{db: db}, nil
}

func (d *SQLDatabase) GetDB() *sql.DB {
	return d.db
}

func (d *SQLDatabase) Close() {
	d.db.Close()
}

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}
