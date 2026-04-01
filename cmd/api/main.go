package main

import (
	"fmt"
	"log"

	"github.com/MdRasB/LogLine/internal/server"
)

func main() {
	fmt.Println("Starting LogLine server...")

	srv := server.NewServer(":8080")
	fmt.Println("Starting the server on :8080")

	err := srv.Start()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
