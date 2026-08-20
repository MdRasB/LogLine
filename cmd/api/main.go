package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MdRasB/LogLine/internal/config"
	"github.com/MdRasB/LogLine/internal/server"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func run() error {
	fmt.Println("Starting LogLine server...")
	startedAt := time.Now()

	cfg, err := config.Load()
	if err != nil {
		log.Println("error in environment configuration")
		return err
	}

	srv, err := server.NewServer(startedAt, cfg)
	if err != nil {
		return err
	}
	fmt.Printf("Starting the server on %v\n", cfg.Port)

	errChan := make(chan error, 1)

	go func() {
		if err := srv.Start(); err != nil {
			errChan <- err
		}
	}()

	stopChan := make(chan os.Signal, 1)

	signal.Notify(stopChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-stopChan:
		{
			fmt.Printf("\nReceived %v signal\n", sig)
			signal.Stop(stopChan)

			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			fmt.Println("Begining graceful shutdown...")
			err := srv.Shutdown(ctx)
			if err != nil {
				return fmt.Errorf("graceful shutdown failed : %w", err)
			}
		}

	case errSig := <-errChan:
		{
			return fmt.Errorf("server failed to start: %w", errSig)
		}

	}

	fmt.Println("The server has been shutdown successfully!")

	return nil
}
