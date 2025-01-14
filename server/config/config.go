package config

import (
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	RedisHost    string
	RedisPort    string
	JwtSecret    string
	ServerPort   string
	ServerHost   string
	KafkaHost    string
	KafkaPort    string
	KafkaTopic   string
	ServerID     string
	BucketName   string
	ObjectKey    string
	AccessKey    string
	SecretAccess string
	Region       string
}

var Config AppConfig

func init() {
	_ = godotenv.Load(".env")
	// Initialize Configuration with environment variables and default values
	Config = AppConfig{
		JwtSecret:    getEnv("JWT_SECRET", "KEY"),
		ServerPort:   getEnv("SERVER_PORT_SERVER", "8082"),
		ServerHost:   getEnv("SERVER_HOST_SERVER", "0.0.0.0"),
		RedisHost:    getEnv("REDIS_HOST", "redis"),
		RedisPort:    getEnv("REDIS_PORT", "6379"),
		KafkaHost:    getEnv("KAFKA_HOST", "kafka"),
		KafkaPort:    getEnv("KAFKA_PORT", "9092"),
		KafkaTopic:   getEnv("KAFKA_TOPIC", "message"),
		ServerID:     uuid.New().String(),
		BucketName:   getEnv("BUCKET_NAME", "shubhamgochat"),
		ObjectKey:    getEnv("OBJECT_KEY", "AbuseMasker/Badwords.txt"),
		AccessKey:    getEnv("ACCESS_KEY", ""),
		SecretAccess: getEnv("SECRET_ACCESS_KEY", ""),
		Region:       getEnv("REGION", ""),
	}
}

// Helper function to fetch environment variables with a default value
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
