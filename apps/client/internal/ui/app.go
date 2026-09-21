package ui

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/Yogesh-Kumar-Mallik-dev/campus-os/apps/client/internal/api"
	"github.com/Yogesh-Kumar-Mallik-dev/campus-os/apps/client/internal/auth"
)

// Screen enumerates top-level application views.
type Screen int

const (
	ScreenLogin Screen = iota
	ScreenOnboarding
	ScreenDashboard
)

// Application encapsulates the Fyne native UI runtime and screen routing.
type Application struct {
	FyneApp       fyne.App
	Window        fyne.Window
	APIClient     *api.Client
	SessionStore  auth.SessionStore
	ActiveSession *auth.AuthSession
	ActiveScreen  Screen
	rootContainer *fyne.Container
}

// BLOCK_UI_APP_NEW_001
// Purpose: Constructs the Campus OS native client application window and sets smart session routing.
func NewApplication(apiClient *api.Client) *Application {
	a := app.NewWithID("internal.campus-os.client")
	a.Settings().SetTheme(NewCampusTheme())

	w := a.NewWindow("Campus OS — Institutional Operating System")
	store := auth.NewKeyringSessionStore("")

	appInst := &Application{
		FyneApp:       a,
		Window:        w,
		APIClient:     apiClient,
		SessionStore:  store,
		rootContainer: container.NewStack(),
	}

	return appInst
}

// NavigateToLogin switches window content to the institutional login screen.
func (a *Application) NavigateToLogin() {
	a.ActiveScreen = ScreenLogin
	a.Window.SetTitle("Campus OS — Institutional Sign In")

	topBar := NewTopBar("BBDIT CAMPUS OS", "Sign In")
	loginView := NewLoginView(
		a.APIClient,
		a.Window,
		a.SessionStore,
		func(session *auth.AuthSession) {
			a.NavigateToDashboard(session)
		},
		func() {
			a.NavigateToOnboarding()
		},
	)

	body := NewResponsiveCardContainer(loginView.CanvasObject(), 560)
	content := container.NewBorder(topBar, nil, nil, nil, body)

	a.rootContainer.Objects = []fyne.CanvasObject{content}
	a.rootContainer.Refresh()
}

// NavigateToOnboarding switches window content to the admission / genesis QR activation wizard.
func (a *Application) NavigateToOnboarding() {
	a.ActiveScreen = ScreenOnboarding
	a.Window.SetTitle("Campus OS — Institutional Account Activation")

	topBar := NewTopBar("BBDIT CAMPUS OS", "Account Activation")
	wizard := NewOnboardingWizard(a.APIClient, a.Window, func() {
		sess := &auth.AuthSession{
			Username: "active.scholar",
			FullName: "Verified Scholar",
			RoleCode: "STUDENT",
		}
		a.NavigateToDashboard(sess)
	})
	wizard.onNavigateToLogin = func() {
		a.NavigateToLogin()
	}

	body := NewResponsiveCardContainer(wizard.CanvasObject(), 640)
	content := container.NewBorder(topBar, nil, nil, nil, body)

	a.rootContainer.Objects = []fyne.CanvasObject{content}
	a.rootContainer.Refresh()
}

// NavigateToDashboard switches window content to the role-based dashboard shell.
func (a *Application) NavigateToDashboard(session *auth.AuthSession) {
	a.ActiveScreen = ScreenDashboard
	a.ActiveSession = session
	a.Window.SetTitle("Campus OS — Executive Portal")

	dashboard := NewDashboardView(
		session,
		a.APIClient,
		a.Window,
		a.SessionStore,
		func() {
			a.NavigateToLogin()
		},
	)

	a.rootContainer.Objects = []fyne.CanvasObject{dashboard.CanvasObject()}
	a.rootContainer.Refresh()
}

// BuildLayout initializes the root view based on smart session discovery.
func (a *Application) BuildLayout() fyne.CanvasObject {
	// Smart launch route: check native OS Keyring
	if a.SessionStore != nil {
		if sess, err := a.SessionStore.Load(); err == nil && sess != nil && sess.AccessToken != "" {
			a.NavigateToDashboard(sess)
			return a.rootContainer
		}
	}

	// Clean launch: default to Login view
	a.NavigateToLogin()
	return a.rootContainer
}

// Run launches the native Fyne window loop.
func (a *Application) Run() {
	a.Window.Resize(fyne.NewSize(1024, 720))
	content := a.BuildLayout()
	a.Window.SetContent(content)
	a.Window.CenterOnScreen()

	// Proactively trigger layout refreshes to eliminate Wayland/Hyprland first-frame stalls and scale lag
	go func() {
		for _, delay := range []time.Duration{40 * time.Millisecond, 120 * time.Millisecond, 250 * time.Millisecond, 450 * time.Millisecond} {
			time.Sleep(delay)
			fyne.Do(func() {
				if a.Window != nil && a.Window.Canvas() != nil {
					if c := a.Window.Content(); c != nil {
						c.Refresh()
						a.Window.Canvas().Refresh(c)
					}
				}
			})
		}
	}()

	a.Window.ShowAndRun()
}
