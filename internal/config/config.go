// Package config creates the server's configuration for http ports and databse ports
package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port  string
	DBURL string

	ReqPerSec float64
	Burst     int

	Version string

	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Printf(".env file not found")
	}

	reqPerSec := getEnv("REQLIMIT", "5")
	burst := getEnv("BURST", "10")

	Port := getPort("PORT", ":8079")
	DBURL := getEnv("DB_URL", "")

	Version := getEnv("VERSION", "v0.0.0")

	ReqPerSec, err := strconv.ParseFloat(reqPerSec, 64)
	if err != nil {
		log.Printf("Failed to convert string: %v", err)
		return nil, err
	}

	Burst, err := strconv.Atoi(burst)
	if err != nil {
		log.Printf("Failed to convert string: %v", err)
		return nil, err
	}

	if Burst < 0 || ReqPerSec < 0 {
		log.Printf("Invalid rate limit. Must be greater or equal to 0.")
		return nil, err
	}

	readTimeout := getEnv("READ_TIMEOUT", "10")
	ReadTimeout, err := time.ParseDuration(readTimeout)
	if err != nil {
		log.Printf("Failed to parse the read time duration: %v", err)
		return nil, err
	}

	readHeaderTimeout := getEnv("READ_HEADER_TIMEOUT", "5")
	ReadHeaderTimeout, err := time.ParseDuration(readHeaderTimeout)
	if err != nil {
		log.Printf("Failed to parse the read header time duration: %v", err)
		return nil, err
	}

	writeTimeout := getEnv("WRITE_TIMEOUT", "15")
	WriteTimeout, err := time.ParseDuration(writeTimeout)
	if err != nil {
		log.Printf("Failed to parse the write time duration: %v", err)
		return nil, err
	}

	idleTimeout := getEnv("IDLE_TIMEOUT", "60")
	IdleTimeout, err := time.ParseDuration(idleTimeout)
	if err != nil {
		log.Printf("Failed to parse the idle time duration: %v", err)
		return nil, err
	}

	if ReadTimeout <= 0 || ReadHeaderTimeout <= 0 || WriteTimeout <= 0 || IdleTimeout <= 0 {
		log.Printf("Invalid request timeout duration. Must be a positive number")
		return nil, err
	}

	cfg := &Config{
		Port:              Port,
		DBURL:             DBURL,
		ReqPerSec:         ReqPerSec,
		Burst:             Burst,
		Version:           Version,
		ReadTimeout:       ReadTimeout,
		ReadHeaderTimeout: ReadHeaderTimeout,
		WriteTimeout:      WriteTimeout,
		IdleTimeout:       IdleTimeout,
	}

	if cfg.DBURL == "" {
		log.Println("DB_URL must be set in environment")
		return nil, err
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return fallback
}

func getPort(key, fallback string) string {
	val := os.Getenv(key)

	if len(val) < 3 || val[0] != ':' {
		log.Println("The port is invalid")
		return fallback
	}

	numVal := val[1:]
	_, err := strconv.Atoi(numVal)
	if err != nil {
		log.Println("The port is invalid")
		return fallback
	}

	return val
}
