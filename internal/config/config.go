package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port  string
	DBURL string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found (this is OK if running in production)")
	}

	cfg := &Config{
		Port:  getEnv("PORT", ":8080"),
		DBURL: getEnv("DB_URL", ""), 
	}

	if cfg.DBURL == "" {
		log.Fatal("DB_URL must be set in environment")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		fmt.Printf("Port is being used: %v\n", val)
		return val
	}

	fmt.Printf("fallback value is being used...%v\n", fallback)
	return fallback
}
