package config

import (
	"os"
)

type AppConfig struct {
	RedisHost        string
	RedisPort        string
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDatabase string
	JwtSecret        string
	ServerPort       string
	KafkaHost        string
	KafkaPort        string
	KafkaTopic       string
	ServerID         string
	AwsRegion        string
}

var Config AppConfig

func init() {
	// Initialize Configuration
	Config = AppConfig{
		JwtSecret:        os.Getenv("JWT_SECRET"),
		ServerPort:       os.Getenv("SERVER_PORT"),
		RedisHost:        os.Getenv("REDIS_HOST"),
		RedisPort:        os.Getenv("REDIS_PORT"),
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDatabase: os.Getenv("POSTGRES_DATABASE"),
		KafkaHost:        os.Getenv("KAFKA_HOST"),
		KafkaPort:        os.Getenv("KAFKA_PORT"),
		KafkaTopic:       os.Getenv("KAFKA_TOPIC"),
		ServerID:         "12345",
		AwsRegion:        os.Getenv("AWS_REGION"),
	}

	// if Config.RedisHost == "" || Config.PostgresHost == "" {
	// 	log.Fatal("Environment variables not set")
	// }
}

// to explicity define values for variables if value not set in .env

func getEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
