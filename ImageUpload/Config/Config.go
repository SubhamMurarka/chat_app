package Config

import "os"

type AppConfig struct {
	AppPort    string
	AppHost    string
	RedisHost  string
	RedisPort  string
	BucketName string
	ObjectKey  string
}

var Conf *AppConfig

func init() {
	Conf = &AppConfig{
		AppPort:    "8089",
		AppHost:    "0.0.0.0",
		BucketName: os.Getenv("BUCKET_NAME"),
		ObjectKey:  os.Getenv("OBJECT_KEY"),
		RedisHost:  "localhost",
		RedisPort:  "6379",
	}
}
