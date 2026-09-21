package ui

import (
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
	topBar := NewTopBar("CAMPUS OS", "Account Activation")

	header := NewPageHeader(
		"Student & Scholar Onboarding",
		"Verify sealed admission QR, validate hardware SIM, and generate your Digital Campus Pass",
		NewStatusPill("SECURE ONBOARDING", PillInfo),
	)

	wizard := NewOnboardingWizard(a.APIClient, a.Window, func() {
		// Callback on onboarding completion
	})

	contentBox := container.NewVBox(
		header,
		wizard.CanvasObject(),
	)

	// Responsive card container adapting fluidly across desktop & mobile screens
	responsiveBody := NewResponsiveCardContainer(contentBox, 660)

	return container.NewBorder(
		topBar,
		nil, nil, nil,
		responsiveBody,
	)
}

// Run launches the native Fyne window loop.
func (a *Application) Run() {
	a.Window.SetContent(a.BuildLayout())
	a.Window.Resize(fyne.NewSize(820, 680))
	a.Window.CenterOnScreen()
	a.Window.ShowAndRun()
}
