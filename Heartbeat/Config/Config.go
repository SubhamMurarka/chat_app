package Config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	RedisHost string
	RedisPort string
}

var Conf *AppConfig

func init() {
	_ = godotenv.Load(".env")
	Conf = &AppConfig{
		RedisHost: getEnv("REDIS_HOST", "redis"), // Default to "redis" if not set
		RedisPort: getEnv("REDIS_PORT", "6379"),  // Default to "6379" if not set
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
