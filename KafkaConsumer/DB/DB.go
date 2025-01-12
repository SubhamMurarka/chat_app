package DB

import (
	"context"
	"log"

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
	pgxConfig.MaxConns = 20

	pgxpool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		log.Fatalf("Error getting config object : %v", err)
		return nil
	}

	return pgxpool
}

func NewSQLDatabase() (*SQLDatabase, error) {
	shards := make([]*pgxpool.Pool, 2)

	pgPool0 := GetConfig("postgres://root:password@localhost:5433/msg_shard0?sslmode=disable")
	pgPool1 := GetConfig("postgres://root:password@localhost:5434/msg_shard1?sslmode=disable")

	shards = append(shards, pgPool0, pgPool1)

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
	ctx := context.Background()
	tx, err := d.db[shardIndex].Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Prepare the statement for bulk insertion
	query := "INSERT INTO message (id, content, user_id, channel_id, message_type, media_id) VALUES ($1, $2, $3, $4, $5, $6)"

	// Create a batch to hold all insert operations
	pgbatch := &pgx.Batch{}
	for _, message := range batch {
		pgbatch.Queue(query, message.ID, message.Content, message.UserID, message.ChannelID, message.MessageType, message.MediaID)
	}

	// Send the batch to the database
	br := tx.SendBatch(ctx, pgbatch)
	defer br.Close()

	// Check for errors in the batch execution
	_, err = br.Exec()
	if err != nil {
		return err
	}

	// Commit the transaction
	return tx.Commit(ctx)
}
