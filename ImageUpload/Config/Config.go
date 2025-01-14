package Config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	AppPort      string
	AppHost      string
	RedisHost    string
	RedisPort    string
	BucketName   string
	ObjectKey    string
	Secret       string
	AccessKey    string
	SecretAccess string
	Region       string
}

var Conf *AppConfig

func init() {
	_ = godotenv.Load(".env")
	Conf = &AppConfig{
		AppPort:      getEnv("SERVER_PORT_IMAGE", "8083"),
		AppHost:      getEnv("SERVER_HOST_IMAGE", "0.0.0.0"),
		BucketName:   getEnv("BUCKET_NAME_IMAGE", "shubhamgochat"),
		ObjectKey:    getEnv("OBJECT_KEY_IMAGE", "ChatImage"),
		RedisHost:    getEnv("REDIS_HOST", "redis"),
		RedisPort:    getEnv("REDIS_PORT", "6379"),
		Secret:       getEnv("SECRET", "KEY"),
		AccessKey:    getEnv("ACCESS_KEY", ""),
		SecretAccess: getEnv("SECRET_ACCESS_KEY", ""),
		Region:       getEnv("REGION", ""),
	}
}

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
