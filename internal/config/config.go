// Package config creates the server's configuration for http ports and databse ports
package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	DBURL     string
	ReqPerSec float64
	Burst     int
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found")
	}

	reqPerSec := getEnv("REQLIMIT", "5")
	burst := getEnv("BURST", "10")

	Port := getEnv("PORT", ":8079")
	DBURL := getEnv("DB_URL", "")

	ReqPerSec, err := strconv.ParseFloat(reqPerSec, 64)
	if err != nil {
		log.Fatalf("Failed to convert string: %v", err)
	}

	Burst, err := strconv.Atoi(burst)
	if err != nil {
		log.Fatalf("Failed to convert string: %v", err)
	}

	if Burst < 0 || ReqPerSec < 0 {
		log.Fatalf("Invalid rate limit. Must be greater or equal to 0.")
	}

	cfg := &Config{
		Port:      Port,
		DBURL:     DBURL,
		ReqPerSec: ReqPerSec,
		Burst:     Burst,
	}

	if cfg.DBURL == "" {
		log.Fatal("DB_URL must be set in environment")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return fallback
}
