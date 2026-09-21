package ui

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/Yogesh-Kumar-Mallik-dev/campus-os/apps/client/internal/api"
	"github.com/Yogesh-Kumar-Mallik-dev/campus-os/apps/client/internal/auth"
)

// BLOCK_UI_DASHBOARD_TEST_001
// Purpose: Verifies DashboardView layout rendering, section switching, and logout trigger.
func TestDashboardView_SectionsAndLogout(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()
	testWindow := test.NewWindow(nil)

	client := api.NewClient("http://localhost:8080")
	tempDir := t.TempDir()
	store := auth.NewKeyringSessionStore(tempDir)

	session := &auth.AuthSession{
		AccessToken:  "jwt_access_mock",
		RefreshToken: "jwt_refresh_mock",
		UserID:       "u-chairperson",
		Username:     "chairperson.2024",
		FullName:     "Yogesh Kumar Mallik",
		RoleCode:     "SUPER_ADMIN",
	}

	logoutTriggered := false

	dashboard := NewDashboardView(
		session,
		client,
		testWindow,
		store,
		func() {
			logoutTriggered = true
		},
	)

	if dashboard.CanvasObject() == nil {
		t.Fatal("expected non-nil canvas object from DashboardView")
	}

	if dashboard.ActiveSection() != SectionOverview {
		t.Errorf("expected initial section to be SectionOverview, got %v", dashboard.ActiveSection())
	}

	// Test section transitions across all 5 navigation destinations
	sections := []DashboardSection{
		SectionAcademicStructure,
		SectionGovernance,
		SectionAuditLedger,
		SectionExecutiveCredentialing,
		SectionOverview,
	}

	for _, s := range sections {
		dashboard.SetSection(s)
		if dashboard.ActiveSection() != s {
			t.Errorf("expected active section %v, got %v", s, dashboard.ActiveSection())
		}
	}

	// Test Logout action
	dashboard.handleLogout()
	if !logoutTriggered {
		t.Errorf("expected logout callback to be triggered")
	}
}

// BLOCK_UI_DASHBOARD_TEST_002
// Purpose: Verifies MetricCard rendering.
func TestMetricCard_Render(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()

	card := NewMetricCard("Total Scholars", "1,420", "Active Enrollment", LucideUserCheck, BadgeSuccess)
	if card == nil {
		t.Fatal("expected non-nil metric card")
	}
}
