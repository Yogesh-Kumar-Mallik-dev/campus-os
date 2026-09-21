package ui

import (
	"testing"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"github.com/Yogesh-Kumar-Mallik-dev/campus-os/apps/client/internal/api"
	"github.com/Yogesh-Kumar-Mallik-dev/campus-os/apps/client/internal/auth"
)

// BLOCK_UI_TEST_001
// Purpose: Verifies headless construction of Fyne UI layout and clean launch into Login view.
func TestAppBuildLayout(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()

	testWindow := test.NewWindow(nil)
	tempDir := t.TempDir()
	store := auth.NewKeyringSessionStore(tempDir)

	clientApp := &Application{
		FyneApp:       testApp,
		Window:        testWindow,
		APIClient:     api.NewClient("http://localhost:8080"),
		SessionStore:  store,
		rootContainer: container.NewStack(),
	}

	layout := clientApp.BuildLayout()
	if layout == nil {
		t.Fatal("BLOCK_UI_TEST_001: expected non-nil layout canvas object")
	}
	if clientApp.ActiveScreen != ScreenLogin {
		t.Errorf("expected clean launch to route to ScreenLogin, got %v", clientApp.ActiveScreen)
	}
}

// BLOCK_UI_TEST_002
// Purpose: Verifies screen routing transitions between Login, Onboarding, and Dashboard.
func TestAppScreenRouting(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()

	testWindow := test.NewWindow(nil)
	tempDir := t.TempDir()
	store := auth.NewKeyringSessionStore(tempDir)

	clientApp := &Application{
		FyneApp:       testApp,
		Window:        testWindow,
		APIClient:     api.NewClient("http://localhost:8080"),
		SessionStore:  store,
		rootContainer: container.NewStack(),
	}

	// 1. Navigate to Onboarding
	clientApp.NavigateToOnboarding()
	if clientApp.ActiveScreen != ScreenOnboarding {
		t.Errorf("expected ScreenOnboarding, got %v", clientApp.ActiveScreen)
	}

	// 2. Navigate to Login
	clientApp.NavigateToLogin()
	if clientApp.ActiveScreen != ScreenLogin {
		t.Errorf("expected ScreenLogin, got %v", clientApp.ActiveScreen)
	}

	// 3. Navigate to Dashboard with session
	sess := &auth.AuthSession{
		AccessToken: "test_token",
		Username:    "chairperson.2024",
		RoleCode:    "SUPER_ADMIN",
	}
	clientApp.NavigateToDashboard(sess)
	if clientApp.ActiveScreen != ScreenDashboard {
		t.Errorf("expected ScreenDashboard, got %v", clientApp.ActiveScreen)
	}
	if clientApp.ActiveSession == nil || clientApp.ActiveSession.RoleCode != "SUPER_ADMIN" {
		t.Errorf("unexpected active session: %+v", clientApp.ActiveSession)
	}

	// 4. Test smart launch with pre-existing session in store
	_ = store.Save(sess)
	clientApp2 := &Application{
		FyneApp:       testApp,
		Window:        testWindow,
		APIClient:     api.NewClient("http://localhost:8080"),
		SessionStore:  store,
		rootContainer: container.NewStack(),
	}
	clientApp2.BuildLayout()
	if clientApp2.ActiveScreen != ScreenDashboard {
		t.Errorf("expected auto-login to ScreenDashboard when session exists, got %v", clientApp2.ActiveScreen)
	}
}
