package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/SubhamMurarka/chat_app/kafkaConsumer/Config"
	"github.com/SubhamMurarka/chat_app/kafkaConsumer/Consumer"
	"github.com/SubhamMurarka/chat_app/kafkaConsumer/DB"
	"github.com/SubhamMurarka/chat_app/kafkaConsumer/MessageIngest"
	models "github.com/SubhamMurarka/chat_app/kafkaConsumer/Models"
)

func main() {
	dbConn, err := DB.NewSQLDatabase()
	if err != nil {
		log.Fatal(err)
		return
	}

	consumer := Consumer.StartConsumer()
	defer consumer.Close()

	batcher := MessageIngest.NewMessageBatcher(dbConn, 100, 10*time.Second)

	partitionConsumer, err := consumer.ConsumePartition(Config.Config.KafkaTopic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("error starting partition consumer: %w", err)
		return
	}
	defer partitionConsumer.Close()

	// Signal handling for graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-stopChan:
			log.Println("Received shutdown signal, stopping consumer.")
			return
		case message := <-partitionConsumer.Messages():
			var msg models.Message
			json.Unmarshal(message.Value, &msg)
			batcher.AddMessage(msg)
		case err := <-partitionConsumer.Errors():
			log.Printf("Error consuming messages: %v", err)
			return
		}
	}
}
