package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Yogesh-Kumar-Mallik-dev/campus-os/backend/internal/config"
	transportGrpc "github.com/Yogesh-Kumar-Mallik-dev/campus-os/backend/internal/transport/grpc"
	transportHttp "github.com/Yogesh-Kumar-Mallik-dev/campus-os/backend/internal/transport/http"
	"github.com/skip2/go-qrcode"
	campusv1 "github.com/Yogesh-Kumar-Mallik-dev/campus-os/backend/pkg/proto/campus/v1"
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

	// Check for Genesis CLI bootstrap subcommand
	if len(os.Args) > 1 && os.Args[1] == "bootstrap-superadmin" {
		runBootstrapCLI(cfg)
		return
	}

	log.Printf("BLOCK_MAIN_ENTRYPOINT_001: Initializing Campus OS Backend v%s (Env: %s)...", Version, cfg.Environment)
	log.Printf("BLOCK_MAIN_ENTRYPOINT_001: DB Unix Socket Target: %s", cfg.DBSocketPath)

	// Attempt connecting to persistence layer over UDS
	var authClient campusv1.AuthServiceClient
	grpcClient, err := transportGrpc.NewPersistenceClient(cfg.DBSocketPath)
	if err == nil {
		authClient = grpcClient.AuthService()
		defer grpcClient.Close()
	} else {
		log.Printf("BLOCK_MAIN_WARN_001: Persistence client link deferred (socket not yet active): %v", err)
	}

	router := transportHttp.NewRouter(transportHttp.RouterConfig{
		Version:    Version,
		AuthClient: authClient,
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

// BLOCK_MAIN_BOOTSTRAP_CLI_001
// Purpose: Implements the bare-metal server CLI command for Genesis Super Admin (Chairperson) bootstrapping.
func runBootstrapCLI(cfg *config.Config) {
	fs := flag.NewFlagSet("bootstrap-superadmin", flag.ExitOnError)
	name := fs.String("name", "Chairperson", "Full legal name of the institutional Chairperson")
	email := fs.String("email", "chairperson@campus.edu", "Official email address of the Chairperson")
	phone := fs.String("phone", "+919876543210", "Registered mobile phone number (used for SIM binding & OTP)")
	outFile := fs.String("out", "genesis_docket.png", "Path to export high-resolution QR docket image file for physical printing")
	blank := fs.Bool("blank", false, "Reset the Super Admin seat to a completely blank, unactivated state for testing")

	_ = fs.Parse(os.Args[2:])

	fmt.Println("========================================================================")
	fmt.Println("  CAMPUS OS GENESIS BOOTSTRAP: SUPER ADMIN (CHAIRPERSON) PROVISIONING   ")
	fmt.Println("========================================================================")
	fmt.Printf("  Target Server Socket : %s\n", cfg.DBSocketPath)
	fmt.Printf("  Chairperson Name     : %s\n", *name)
	fmt.Printf("  Chairperson Email    : %s\n", *email)
	fmt.Printf("  Registered Phone     : %s\n", *phone)
	fmt.Printf("  Blank Sheet Reset    : %t\n", *blank)
	fmt.Println("------------------------------------------------------------------------")

	grpcClient, err := transportGrpc.NewPersistenceClient(cfg.DBSocketPath)
	if err != nil {
		fmt.Printf("Error: Failed to dial DB layer over %s: %v\n", cfg.DBSocketPath, err)
		os.Exit(1)
	}
	defer grpcClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := grpcClient.AuthService().BootstrapSuperAdmin(ctx, &campusv1.BootstrapSuperAdminRequest{
		ChairpersonName:  *name,
		ChairpersonEmail: *email,
		ChairpersonPhone: *phone,
		Blank:            *blank,
	})
	if err != nil {
		fmt.Printf("Bootstrap Error: %v\n", err)
		os.Exit(1)
	}

	// Export printable high-resolution QR docket image
	if *outFile != "" {
		if err := qrcode.WriteFile(resp.ClaimToken, qrcode.Medium, 512, *outFile); err == nil {
			fmt.Printf("  Exported Docket File   : %s (High-Resolution 512x512 PNG)\n", *outFile)
		} else {
			fmt.Printf("  Warning: Failed to export docket image to %s: %v\n", *outFile, err)
		}
	}

	fmt.Println("\n  SCAN THIS SEALED CLAIM QR CODE USING CAMPUS OS CLIENT:")
	fmt.Println()
	if qr, err := qrcode.New(resp.ClaimToken, qrcode.Medium); err == nil {
		fmt.Println(qr.ToSmallString(false))
	} else if resp.AsciiQr != "" {
		fmt.Println(resp.AsciiQr)
	}
	fmt.Println()
	fmt.Printf("  Single-Use Claim Token : %s\n", resp.ClaimToken)
	fmt.Printf("  Expires At (Unix)      : %d\n", resp.ExpiresAtUnix)
	fmt.Println("========================================================================")
	fmt.Println("  INSTRUCTIONS:")
	fmt.Println("  1. Open Campus OS Native Client on the Chairperson's registered device.")
	fmt.Println("  2. Point camera at the QR code above or upload the exported genesis_docket.png.")
	fmt.Println("  3. Verify cellular SIM binding and set your master administrator password.")
	fmt.Println("========================================================================")
}
