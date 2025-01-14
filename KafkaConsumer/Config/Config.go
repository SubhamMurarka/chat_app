package Config

import (
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	KafkaHost        string
	KafkaPort        string
	KafkaTopic       string
	ServerID         string
	PostgresHost1    string
	PostgresHost2    string
	PostgresUser     string
	PostgresPort     string
	PostgresPassword string
}

var Config AppConfig

func init() {
	_ = godotenv.Load(".env")
	// Initialize Configuration
	Config = AppConfig{
		KafkaHost:        getEnv("KAFKA_HOST", "kafka"),
		KafkaPort:        getEnv("KAFKA_PORT", "9092"),
		KafkaTopic:       getEnv("KAFKA_TOPIC", "message"),
		ServerID:         uuid.New().String(),
		PostgresHost1:    getEnv("POSTGRES_HOST_SHARD0", "postgres_msg_shard0"),
		PostgresHost2:    getEnv("POSTGRES_HOST_SHARD1", "postgres_msg_shard1"),
		PostgresUser:     getEnv("POSTGRES_USER", "root"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "password"), // Always generates a new unique ID
	}
}

// to explicity define values for variables if value not set in .env

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		if value == "" {
			log.Printf("Environment variable %s is empty; using default value: %s", key, defaultValue)
			return defaultValue
		}
		log.Printf("Found environment variable %s with value: %s", key, value)
		return value
	}
	log.Printf("Environment variable %s not found; using default value: %s", key, defaultValue)
	return defaultValue
}
