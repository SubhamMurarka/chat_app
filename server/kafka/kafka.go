package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/SubhamMurarka/chat_app/server/config"
	"github.com/SubhamMurarka/chat_app/server/models"
	"github.com/segmentio/kafka-go"
)

type KafkaConfig struct {
	Host  string
	Port  string
	Topic string
}

var kafkaConfig KafkaConfig

func init() {
	kafkaConfig = KafkaConfig{
		Host:  config.Config.KafkaHost,
		Port:  config.Config.KafkaPort,
		Topic: config.Config.KafkaTopic,
	}
}

// var producer sarama.SyncProducer

func ProduceToKafka(message models.Message) {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaConfig.Host + ":" + kafkaConfig.Port},
		Topic:   kafkaConfig.Topic,
	})

	defer w.Close()

	//serialize the message
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Println("Error marshalling message:", err)
		return
	}

	err = w.WriteMessages(context.Background(), kafka.Message{
		Value: messageBytes,
	})

	if err != nil {
		log.Println("Error publishing message to kafka:", err)
		return
	}

}
