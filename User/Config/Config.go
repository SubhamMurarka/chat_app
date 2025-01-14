package Config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	PostgresHost1    string
	PostgresHost2    string
	PostgresHost3    string
	PostgresUser     string
	PostgresPort     string
	PostgresPassword string
	RedisPort        string
	RedisHost        string
	ServerPort       string
	ServerHost       string
	JwtSecret        string
}

var Config *AppConfig

func init() {
	_ = godotenv.Load(".env")
	// Initialize Configuration
	Config = &AppConfig{
		JwtSecret:        getEnv("JWT_SECRET", "KEY"),
		PostgresHost1:    getEnv("POSTGRES_HOST_SHARD0", "postgres_msg_shard0"),
		PostgresHost2:    getEnv("POSTGRES_HOST_SHARD1", "postgres_msg_shard1"),
		PostgresHost3:    getEnv("POSTGRES_HOST_USER", "postgres_user"),
		PostgresUser:     getEnv("POSTGRES_USER", "root"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "password"),
		RedisHost:        getEnv("REDIS_HOST", "redis"),
		RedisPort:        getEnv("REDIS_PORT", "6379"),
		ServerPort:       getEnv("SERVER_PORT_USER", "8081"),
		ServerHost:       getEnv("SERVER_HOST_USER", "0.0.0.0"),
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
