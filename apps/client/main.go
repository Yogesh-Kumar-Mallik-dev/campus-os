package main

import (
	"os"

	"github.com/Yogesh-Kumar-Mallik-dev/campus-os/apps/client/internal/api"
	"github.com/Yogesh-Kumar-Mallik-dev/campus-os/apps/client/internal/ui"
)

// BLOCK_CLIENT_MAIN_001
// Purpose: Baseline entrypoint for Campus OS native cross-platform client using Fyne.
func main() {
	backendURL := os.Getenv("CAMPUS_BACKEND_URL")
	if backendURL == "" {
		backendURL = "http://localhost:8080"
	}

	apiClient := api.NewClient(backendURL)
	application := ui.NewApplication(apiClient)

	application.Run()
}
