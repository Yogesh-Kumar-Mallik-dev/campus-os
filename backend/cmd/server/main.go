package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// BLOCK_MAIN_ENTRYPOINT_001
// Purpose: Initializes the Campus OS Go Logical Backend, loads config, and orchestrates domain services.
// Inputs:  Environment variables, runtime flags
// Outputs: Graceful server shutdown on SIGINT/SIGTERM
// Errors:  ERR_INIT_FAILURE
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("BLOCK_MAIN_ENTRYPOINT_001: Starting Campus OS Logical Backend...")

	// Listen for shutdown signal
	go func() {
		sig := <-sigChan
		log.Printf("BLOCK_MAIN_ENTRYPOINT_001: Received signal %v, initiating shutdown...", sig)
		cancel()
	}()

	<-ctx.Done()
	fmt.Println("BLOCK_MAIN_ENTRYPOINT_001: Campus OS Logical Backend stopped cleanly.")
}
