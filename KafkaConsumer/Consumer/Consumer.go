package Consumer

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/SubhamMurarka/chat_app/kafkaConsumer/Config"
)

func StartConsumer() sarama.Consumer {
	// Initialize Kafka consumer
	consumer, err := initializeKafkaConsumer()
	if err != nil {
		log.Fatalf("Error Handling Kafka Consumer : %v", err)
		return nil
	}

	return consumer
}

func initializeKafkaConsumer() (sarama.Consumer, error) {
	url := fmt.Sprintf("%s:%s", Config.Config.KafkaHost, Config.Config.KafkaPort)
	brokerUrls := []string{url}
	fmt.Println(url)

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	var consumer sarama.Consumer
	var err error
	for retries := 10; retries > 0; retries-- {
		consumer, err = sarama.NewConsumer(brokerUrls, config)
		if err == nil {
			log.Println("kafka consumer connected successfully")
			return consumer, nil
		}
		log.Printf("Error creating Kafka consumer, retrying: %v", err)
		time.Sleep(10 * time.Second)
	}
	return nil, fmt.Errorf("failed to create Kafka consumer after retries: %w", err)
}

// func processMessages(consumer sarama.Consumer, topic string, esClient *elasticsearch.Client, batchSize int) error {
// 	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
// 	if err != nil {
// 		return fmt.Errorf("error starting partition consumer: %w", err)
// 	}
// 	defer partitionConsumer.Close()

// 	// Signal handling for graceful shutdown
// 	stopChan := make(chan os.Signal, 1)
// 	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

// 	batch := make([]models.ProcessedData, 0, batchSize)

// 	for {
// 		select {
// 		case <-stopChan:
// 			log.Println("Received shutdown signal, stopping consumer.")
// 			return nil
// 		case message := <-partitionConsumer.Messages():
// 			var data models.AnalyticsData
// 			if err := json.Unmarshal(message.Value, &data); err != nil {
// 				log.Printf("Error unmarshaling message: %v", err)
// 				continue
// 			}
// 			fmt.Println("data : ", data)
// 			processed := processAnalyticsData(data)
// 			fmt.Println("processed : ", processed)
// 			batch = append(batch, processed)

// 			if len(batch) >= batchSize {
// 				if err := es.WriteBatchToElasticsearch(batch, esClient); err != nil {
// 					log.Printf("Error writing batch to Elasticsearch: %v", err)
// 				}
// 				batch = batch[:0]
// 			}
// 		case err := <-partitionConsumer.Errors():
// 			log.Printf("Error consuming messages: %v", err)
// 		}
// 	}
// }
