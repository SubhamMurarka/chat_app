package MessageIngest

import (
	"fmt"
	"sync"
	"time"

	"github.com/SubhamMurarka/chat_app/kafkaConsumer/DB"
	models "github.com/SubhamMurarka/chat_app/kafkaConsumer/Models"
)

type MessageBatcher struct {
	Db           *DB.SQLDatabase
	BatchSize    int
	ShardBatches [][]models.Message
	Mu           sync.Mutex
	Ticker       *time.Ticker
}

func NewMessageBatcher(db *DB.SQLDatabase, batchSize int, flushInterval time.Duration) *MessageBatcher {
	shardCount := len(db.GetDB())
	shardBatches := make([][]models.Message, shardCount)
	for i := range shardBatches {
		shardBatches[i] = make([]models.Message, 0)
	}

	mb := &MessageBatcher{
		Db:           db,
		BatchSize:    batchSize,
		ShardBatches: shardBatches,
		Ticker:       time.NewTicker(flushInterval),
	}
	go mb.startBatching()
	return mb
}

func (mb *MessageBatcher) AddMessage(message models.Message) {
	mb.Mu.Lock()
	defer mb.Mu.Unlock()

	shardIndex := int(message.ChannelID % int64(len(mb.ShardBatches)))
	mb.ShardBatches[shardIndex] = append(mb.ShardBatches[shardIndex], message)

	if len(mb.ShardBatches[shardIndex]) >= mb.BatchSize {
		mb.flushBatch(shardIndex)
	}
}

func (mb *MessageBatcher) startBatching() {
	for range mb.Ticker.C {
		mb.Mu.Lock()
		for i := range mb.ShardBatches {
			if len(mb.ShardBatches[i]) > 0 {
				mb.flushBatch(i)
			}
		}
		mb.Mu.Unlock()
	}
}

func (mb *MessageBatcher) flushBatch(shardIndex int) {
	batch := mb.ShardBatches[shardIndex]
	if len(batch) > 0 {
		err := mb.Db.InsertMessagesToShard(int64(shardIndex), batch)
		if err != nil {
			fmt.Printf("Error batching the data : %v", err)
			return
		}
		// Clear the batch after processing
		mb.ShardBatches[shardIndex] = nil
	}
}
