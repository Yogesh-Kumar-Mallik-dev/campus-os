package main

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/Yogesh-Kumar-Mallik-dev/campus-os/apps/client/internal/api"
	"github.com/Yogesh-Kumar-Mallik-dev/campus-os/apps/client/internal/layout"
	"github.com/Yogesh-Kumar-Mallik-dev/campus-os/apps/client/internal/ui"
)

// BLOCK_CLIENT_MAIN_001
// Purpose: Baseline entrypoint for Campus OS native cross-platform client using Fyne.
func main() {
	for _, arg := range os.Args[1:] {
		if arg == "--demo" || arg == "demo" {
			a := app.NewWithID("com.campusos.layout.demo")
			w := a.NewWindow("Campus OS — Fully Adaptive Responsive Layout Engine")
			w.SetContent(layout.BuildDemoScreen(w))
			w.Resize(fyne.NewSize(1024, 720))
			w.CenterOnScreen()
			w.ShowAndRun()
			return
		}
	}

	backendURL := os.Getenv("CAMPUS_BACKEND_URL")
	if backendURL == "" {
		backendURL = "http://localhost:8080"
	}

	apiClient := api.NewClient(backendURL)
	application := ui.NewApplication(apiClient)

	application.Run()
}
