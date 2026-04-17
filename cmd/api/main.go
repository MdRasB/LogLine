package main

import (
	"fmt"
	"log"

	"github.com/MdRasB/LogLine/internal/config"
	"github.com/MdRasB/LogLine/internal/server"
)

func main() {
	fmt.Println("Starting LogLine server...")

	cfg := config.Load()
	port := cfg.Port
	dbStr := cfg.DBURL
	
	srv := server.NewServer(port, dbStr)
	fmt.Printf("Starting the server on %v", port)

	err := srv.Start()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
