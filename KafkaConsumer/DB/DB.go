package DB

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/SubhamMurarka/chat_app/kafkaConsumer/Config"
	models "github.com/SubhamMurarka/chat_app/kafkaConsumer/Models"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SQLDatabase struct {
	db []*pgxpool.Pool
}

func GetConfig(url string) *pgxpool.Pool {
	pgxConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatalf("Error getting config object : %v", err)
		return nil
	}

	pgxConfig.MinConns = 5
	pgxConfig.MaxConns = 48

	pgxpool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		log.Fatalf("Error getting config object : %v", err)
		return nil
	}

	return pgxpool
}

func NewSQLDatabase() (*SQLDatabase, error) {
	shards := make([]*pgxpool.Pool, 2)

	commonstring := fmt.Sprintf("postgres://%s:%s@", Config.Config.PostgresUser, Config.Config.PostgresPassword)
	pgPool0 := GetConfig(commonstring + Config.Config.PostgresHost1 + ":" + Config.Config.PostgresPort + "/msg_shard0?sslmode=disable")
	pgPool1 := GetConfig(commonstring + Config.Config.PostgresHost2 + ":" + Config.Config.PostgresPort + "/msg_shard1?sslmode=disable")
	shards[0] = pgPool0
	shards[1] = pgPool1

	return &SQLDatabase{db: shards}, nil
}

func (d *SQLDatabase) GetDB() []*pgxpool.Pool {
	return d.db
}

func (d *SQLDatabase) Close() {
	d.db[0].Close()
	d.db[1].Close()
}

// func (d *SQLDatabase) InsertMessagesToShard(shardIndex int64, batch []models.Message) error {
// 	query :=
// }

func (d *SQLDatabase) InsertMessagesToShard(shardIndex int64, batch []models.Message) error {
	if len(batch) == 0 {
		return nil
	}

	ctx := context.Background()
	tx, err := d.db[shardIndex].Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := "INSERT INTO messages (id, content, user_id, channel_id, media_id, message_type) VALUES ($1, $2, $3, $4, $5, $6)"

	pgbatch := &pgx.Batch{}
	for _, message := range batch {
		pgbatch.Queue(query, message.ID, message.Content, message.UserID, message.ChannelID, message.MediaID, message.MessageType)
	}

	br := tx.SendBatch(ctx, pgbatch)
	for i := 0; i < len(batch); i++ {
		if _, err := br.Exec(); err != nil {
			br.Close()
			return err
		}
	}
	br.Close()

	return tx.Commit(ctx)
}

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
