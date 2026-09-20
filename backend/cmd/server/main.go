package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Yogesh-Kumar-Mallik-dev/campus-os/backend/internal/config"
	transportHttp "github.com/Yogesh-Kumar-Mallik-dev/campus-os/backend/internal/transport/http"
)

// Version of the Campus OS Backend
const Version = "0.1.0"

// BLOCK_MAIN_ENTRYPOINT_001
// Purpose: Boots the Campus OS Go Logical Backend, binds HTTP routes, and orchestrates lifecycle.
// Inputs:  Environment configuration
// Outputs: Graceful server shutdown on OS signal
// Errors:  Logs and exits on unrecoverable bind errors
func main() {
	cfg := config.Load()

	log.Printf("BLOCK_MAIN_ENTRYPOINT_001: Initializing Campus OS Backend v%s (Env: %s)...", Version, cfg.Environment)
	log.Printf("BLOCK_MAIN_ENTRYPOINT_001: DB Unix Socket Target: %s", cfg.DBSocketPath)

	router := transportHttp.NewRouter(transportHttp.RouterConfig{
		Version: Version,
	})

	server := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	// Server runner goroutine
	go func() {
		log.Printf("BLOCK_MAIN_ENTRYPOINT_001: HTTP API listening on port %s", cfg.HTTPPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("BLOCK_MAIN_ENTRYPOINT_001: Server listen failed: %v", err)
		}
	}()

	// Graceful shutdown listener
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan

	log.Println("BLOCK_MAIN_ENTRYPOINT_001: Shutdown signal received, draining active connections...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("BLOCK_MAIN_ENTRYPOINT_001: Graceful shutdown error: %v", err)
	}

	fmt.Println("BLOCK_MAIN_ENTRYPOINT_001: Campus OS Backend terminated cleanly.")
}
