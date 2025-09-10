package config

import "os"

type Config struct {
	MongoURI   string
	DBName     string
	ServerAddr string
	SecretKey string
}

func Load() *Config {
	return &Config{
		MongoURI:   getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DBName:     getEnv("MONGO_DB", "url_shortener"),
		ServerAddr: getEnv("SERVER_ADDR", ":8080"),
		SecretKey: getEnv("SECRET", "123"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
