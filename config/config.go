package config

type Config struct {
	HTTPPort string

	PostgresHost           string
	PostgresUser           string
	PostgresDatabase       string
	PostgresPassword       string
	PostgresPort           string
	PostgresMaxConnections int32

	RedisAddr     string
	RedisPassword string
	RedisDB       int

	AuthSecretKey string
	SuperAdmin    string
	Client        string
}

func Load() Config {

	var cfg Config

	cfg.HTTPPort = ":4000"

	cfg.PostgresHost = "localhost"
	cfg.PostgresUser = "samandar"
	cfg.PostgresDatabase = "book"
	cfg.PostgresPassword = "samandevop"
	cfg.PostgresPort = "5432"
	cfg.PostgresMaxConnections = 20

	cfg.RedisAddr = "localhost:6379"
	cfg.RedisPassword = ""
	cfg.RedisDB = 0

	cfg.AuthSecretKey = "9K+WgNTglA44Hg=="
	cfg.SuperAdmin = "Super"
	cfg.Client = "Client"

	return cfg
}
