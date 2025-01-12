package Repository

import (
	"context"

	models "github.com/SubhamMurarka/chat_app/User/Models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RepoInterface interface {
	GetNewMessages(ctx context.Context, channelID int64) ([]models.Message, error)
	PaginationMessages(ctx context.Context, channelID int64, lastid uint64) ([]models.Message, error)
}

type repo struct {
	shards []*pgxpool.Pool
}

func NewRepo(shards []*pgxpool.Pool) RepoInterface {
	return &repo{
		shards: shards,
	}
}

func (r *repo) GetNewMessages(ctx context.Context, channelID int64) ([]models.Message, error) {
	shard := r.determineShard(channelID) // Assuming you determine the shard based on channel ID or other criteria

	query := `
			   SELECT * 
			   FROM messages 
			   WHERE channel_id = $1
			   ORDER BY id DESC 
			   LIMIT 20;
					       `

	rows, err := shard.Query(ctx, query, channelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.ChannelID, &msg.UserID, &msg.Content, &msg.MediaID, &msg.MessageType, &msg.Created_At); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *repo) PaginationMessages(ctx context.Context, channelID int64, lastid uint64) ([]models.Message, error) {
	shard := r.determineShard(channelID)

	query := `
				SELECT * 
				FROM messages
				WHERE id < $1 AND channel_id = $2
				ORDER BY id DESC
				LIMIT 20
	`

	rows, err := shard.Query(ctx, query, lastid, channelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.ChannelID, &msg.UserID, &msg.Content, &msg.MediaID, &msg.MessageType, &msg.Created_At); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

// determineShard selects the appropriate shard based on the channelID.
func (r *repo) determineShard(channelID int64) *pgxpool.Pool {
	index := channelID % int64(len(r.shards))
	return r.shards[index]
}
