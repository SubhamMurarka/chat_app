package Config

import (
	"os"

	"github.com/google/uuid"
)

type AppConfig struct {
	ServerPort string
	KafkaHost  string
	KafkaPort  string
	KafkaTopic string
	ServerID   string
}

var Config AppConfig

func init() {
	// Initialize Configuration
	Config = AppConfig{

		ServerPort: os.Getenv("SERVER_PORT"),
		KafkaHost:  os.Getenv("KAFKA_HOST"),
		KafkaPort:  os.Getenv("KAFKA_PORT"),
		KafkaTopic: os.Getenv("KAFKA_TOPIC"),
		ServerID:   uuid.New().String(),
	}

}

// to explicity define values for variables if value not set in .env

func getEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
