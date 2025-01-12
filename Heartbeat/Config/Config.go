package Config

type AppConfig struct {
	RedisHost string
	RedisPort string
}

var Conf *AppConfig

func init() {
	Conf = &AppConfig{
		RedisHost: "localhost",
		RedisPort: "6379",
	}
}
