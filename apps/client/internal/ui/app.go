package ui

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/Yogesh-Kumar-Mallik-dev/campus-os/apps/client/internal/api"
)

// Application encapsulates the Fyne native UI runtime.
type Application struct {
	FyneApp   fyne.App
	Window    fyne.Window
	APIClient *api.Client
}

// BLOCK_UI_APP_NEW_001
// Purpose: Constructs the Campus OS native client application window, applies custom design tokens, and sets responsive geometry.
func NewApplication(apiClient *api.Client) *Application {
	a := app.NewWithID("internal.campus-os.client")
	a.Settings().SetTheme(NewCampusTheme())

	w := a.NewWindow("Campus OS — Institutional Account Activation")

	return &Application{
		FyneApp:   a,
		Window:    w,
		APIClient: apiClient,
	}
}

// BLOCK_UI_APP_BUILD_001
// Purpose: Assembles fully responsive Account Activation layout matching institutional design tokens.
func (a *Application) BuildLayout() fyne.CanvasObject {
	topBar := NewTopBar("BBDIT CAMPUS OS", "Account Activation")

	wizard := NewOnboardingWizard(a.APIClient, a.Window, func() {
		// Callback on onboarding completion
	})

	// Responsive card container adapting fluidly across desktop & mobile screens
	responsiveBody := NewResponsiveCardContainer(wizard.CanvasObject(), 640)

	return container.NewBorder(
		topBar,
		nil, nil, nil,
		responsiveBody,
	)
}

// Run launches the native Fyne window loop.
func (a *Application) Run() {
	a.Window.Resize(fyne.NewSize(820, 680))
	content := a.BuildLayout()
	a.Window.SetContent(content)
	a.Window.CenterOnScreen()

	// Proactively trigger a layout refresh to eliminate Wayland/Hyprland first-frame stalls
	go func() {
		for _, delay := range []time.Duration{30 * time.Millisecond, 120 * time.Millisecond} {
			time.Sleep(delay)
			if a.Window != nil && a.Window.Canvas() != nil {
				if c := a.Window.Content(); c != nil {
					a.Window.Canvas().Refresh(c)
				}
			}
		}
	}()

	a.Window.ShowAndRun()
}
