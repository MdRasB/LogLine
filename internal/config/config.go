package config

import "os"

type Config struct {
	Port  string
	DBURL string
}

func Load() *Config {
	return &Config{
		Port:  getEnv("PORT", ":8080"),
		DBURL: getEnv("DB_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}