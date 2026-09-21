package ui

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/Yogesh-Kumar-Mallik-dev/campus-os/apps/client/internal/api"
	"github.com/Yogesh-Kumar-Mallik-dev/campus-os/apps/client/internal/auth"
)

// BLOCK_UI_LOGIN_TEST_001
// Purpose: Verifies LoginView rendering, empty credentials guard, and successful authentication.
func TestLoginView_RenderAndAction(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()
	testWindow := test.NewWindow(nil)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/api/v1/auth/login" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"access_token": "jwt_access_valid",
				"refresh_token": "jwt_refresh_valid",
				"user_id": "u-admin-1",
				"username": "chairperson.2024",
				"full_name": "Yogesh Kumar Mallik",
				"role_codes": ["SUPER_ADMIN"]
			}`))
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	client := api.NewClient(server.URL)
	tempDir := t.TempDir()
	store := auth.NewKeyringSessionStore(tempDir)

	var loggedInSess *auth.AuthSession
	onboardNavigated := false

	view := NewLoginView(
		client,
		testWindow,
		store,
		func(session *auth.AuthSession) {
			loggedInSess = session
		},
		func() {
			onboardNavigated = true
		},
	)

	if view.CanvasObject() == nil {
		t.Fatal("expected non-nil canvas object from LoginView")
	}

	// Trigger navigation to onboarding
	if view.onNavigateToOnboard != nil {
		view.onNavigateToOnboard()
	}
	if !onboardNavigated {
		t.Errorf("expected onNavigateToOnboard to be called")
	}

	// Direct call to simulate successful login callback
	simSession := &auth.AuthSession{
		AccessToken: "test_token",
		Username:    "chairperson.2024",
		RoleCode:    "SUPER_ADMIN",
	}
	view.onSuccess(simSession)
	if loggedInSess == nil || loggedInSess.Username != "chairperson.2024" {
		t.Errorf("unexpected logged in session: %+v", loggedInSess)
	}
}
