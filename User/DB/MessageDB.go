package db

import (
	"context"
	"fmt"
	"log"

	"github.com/SubhamMurarka/chat_app/User/Config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MsgDatabase struct {
	db []*pgxpool.Pool
}

func GetConfig(url string) *pgxpool.Pool {
	pgxConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatalf("Error getting config object : %v", err)
		return nil
	}

	pgxConfig.MinConns = 5
	pgxConfig.MaxConns = 20

	pgxpool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		log.Fatalf("Error getting config object : %v", err)
		return nil
	}

	return pgxpool
}

func NewMsgDatabase() (*MsgDatabase, error) {
	shards := make([]*pgxpool.Pool, 2)

	commonstring := fmt.Sprintf("postgres://%s:%s@", Config.Config.PostgresUser, Config.Config.PostgresPassword)
	pgPool0 := GetConfig(commonstring + Config.Config.PostgresHost1 + ":" + Config.Config.PostgresPort + "/msg_shard0?sslmode=disable")
	pgPool1 := GetConfig(commonstring + Config.Config.PostgresHost2 + ":" + Config.Config.PostgresPort + "/msg_shard1?sslmode=disable")

	shards[0] = pgPool0
	shards[1] = pgPool1

	return &MsgDatabase{db: shards}, nil
}

func (m *MsgDatabase) GetDB() []*pgxpool.Pool {
	return m.db
}

func (m *MsgDatabase) Close() {
	m.db[0].Close()
	m.db[1].Close()
}
