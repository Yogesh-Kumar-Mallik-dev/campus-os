package ui

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/test"
	"github.com/Yogesh-Kumar-Mallik-dev/campus-os/apps/client/internal/api"
	"github.com/Yogesh-Kumar-Mallik-dev/campus-os/apps/client/internal/auth"
)

// BLOCK_UI_DASHBOARD_TEST_001
// Purpose: Verifies DashboardView layout rendering, section switching, sidebar toggle, and logout trigger.
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

	// Test section transitions across all navigation destinations
	sections := []DashboardSection{
		SectionAcademicStructure,
		SectionGovernance,
		SectionAuditLedger,
		SectionExecutiveCredentialing,
		SectionSettings,
		SectionOverview,
	}

	for _, s := range sections {
		dashboard.SetSection(s)
		if dashboard.ActiveSection() != s {
			t.Errorf("expected active section %v, got %v", s, dashboard.ActiveSection())
		}
	}

	// Test sidebar collapsing into compact icon rail mode
	dashboard.sidebarOpen = false
	dashboard.buildShell()
	if dashboard.CanvasObject() == nil {
		t.Fatal("expected non-nil canvas object when sidebar is collapsed")
	}

	dashboard.sidebarOpen = true
	dashboard.buildShell()
	if dashboard.CanvasObject() == nil {
		t.Fatal("expected non-nil canvas object when sidebar is expanded")
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

// BLOCK_UI_DASHBOARD_TEST_003
// Purpose: Verifies DashboardView responsive reflow across desktop, tablet/tiled, and narrow mobile viewports.
func TestDashboardView_ResponsiveBreakpoints(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()
	testWindow := test.NewWindow(nil)

	client := api.NewClient("http://localhost:8080")
	tempDir := t.TempDir()
	store := auth.NewKeyringSessionStore(tempDir)

	session := &auth.AuthSession{
		AccessToken: "jwt_access_mock",
		RoleCode:    "SUPER_ADMIN",
		FullName:    "Yogesh Kumar Mallik",
	}

	dashboard := NewDashboardView(
		session,
		client,
		testWindow,
		store,
		func() {},
	)

	// Case 1: Wide Desktop (1200px)
	dashboard.handleWindowResize(fyne.NewSize(1200, 800))
	if !dashboard.sidebarOpen {
		t.Errorf("expected sidebar to be open on wide viewport (1200px)")
	}

	// Case 2: Tiled / Half-Screen Window (850px)
	// Must automatically collapse sidebar and switch to compact topbar / 2-col KPI grid
	dashboard.handleWindowResize(fyne.NewSize(850, 720))
	if dashboard.sidebarOpen {
		t.Errorf("expected sidebar to auto-collapse on tiled viewport (850px)")
	}

	// Case 3: Narrow Window (500px)
	// Must retain collapsed sidebar and switch to 1-col KPI grid
	dashboard.handleWindowResize(fyne.NewSize(500, 720))
	if dashboard.sidebarOpen {
		t.Errorf("expected sidebar to remain collapsed on narrow viewport (500px)")
	}

	// Case 4: Maximize back to Fullscreen (1400px)
	// Must automatically expand sidebar back to 240px
	dashboard.handleWindowResize(fyne.NewSize(1400, 900))
	if !dashboard.sidebarOpen {
		t.Errorf("expected sidebar to auto-expand on fullscreen viewport (1400px)")
	}

	// Case 5: User explicitly toggles sidebar closed on desktop
	dashboard.userToggledSidebar = true
	dashboard.sidebarOpen = false
	dashboard.handleWindowResize(fyne.NewSize(1400, 900))
	if dashboard.sidebarOpen {
		t.Errorf("expected user manual sidebar toggle to be respected regardless of viewport width")
	}
}

// BLOCK_UI_DASHBOARD_TEST_004
// Purpose: Verifies AdaptiveGridLayout computes correct columns and heights based on available width.
func TestAdaptiveGridLayout_DynamicColumnsAndMinSize(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()

	item1 := canvas.NewRectangle(color.White)
	item1.SetMinSize(fyne.NewSize(200, 100))
	item2 := canvas.NewRectangle(color.White)
	item2.SetMinSize(fyne.NewSize(200, 100))
	item3 := canvas.NewRectangle(color.White)
	item3.SetMinSize(fyne.NewSize(200, 100))
	item4 := canvas.NewRectangle(color.White)
	item4.SetMinSize(fyne.NewSize(200, 100))

	items := []fyne.CanvasObject{item1, item2, item3, item4}
	layout := NewAdaptiveGridLayout(200, 10)

	// At 850px width, 4 items of 200px + gap should fit on 1 row (4 cols)
	layout.Layout(items, fyne.NewSize(850, 400))
	if item1.Position().Y != item2.Position().Y || item2.Position().Y != item4.Position().Y {
		t.Errorf("expected all 4 items on row 0 at 850px width, got y positions: %v, %v, %v", item1.Position().Y, item2.Position().Y, item4.Position().Y)
	}

	// At 450px width, items should wrap into 2 columns (2 rows)
	layout.Layout(items, fyne.NewSize(450, 400))
	if item1.Position().Y != item2.Position().Y {
		t.Errorf("expected item1 and item2 on row 0, got y1=%v, y2=%v", item1.Position().Y, item2.Position().Y)
	}
	if item3.Position().Y <= item1.Position().Y {
		t.Errorf("expected item3 to wrap to row 1, got y3=%v, y1=%v", item3.Position().Y, item1.Position().Y)
	}

	// At 200px width, items should wrap into 1 column (4 rows)
	layout.Layout(items, fyne.NewSize(200, 600))
	if item2.Position().Y <= item1.Position().Y {
		t.Errorf("expected item2 to wrap to row 1 at 200px, got y1=%v, y2=%v", item1.Position().Y, item2.Position().Y)
	}
	if item3.Position().Y <= item2.Position().Y {
		t.Errorf("expected item3 to wrap to row 2 at 200px, got y2=%v, y3=%v", item2.Position().Y, item3.Position().Y)
	}

	minSize := layout.MinSize(items)
	if minSize.Height <= 100 {
		t.Errorf("expected minSize height to reflect multiple rows when wrapped, got %v", minSize.Height)
	}
}

// BLOCK_UI_DASHBOARD_TEST_005
// Purpose: Verifies FlowLayout wraps items horizontally when width is exceeded.
func TestFlowLayout_Wrapping(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()

	btn1 := canvas.NewRectangle(color.White)
	btn1.SetMinSize(fyne.NewSize(150, 40))
	btn2 := canvas.NewRectangle(color.White)
	btn2.SetMinSize(fyne.NewSize(150, 40))
	btn3 := canvas.NewRectangle(color.White)
	btn3.SetMinSize(fyne.NewSize(150, 40))

	btns := []fyne.CanvasObject{btn1, btn2, btn3}
	flow := NewFlowLayout(10)

	// At 500px, all three 150px buttons fit in 1 row (150*3 + 20 = 470px)
	flow.Layout(btns, fyne.NewSize(500, 100))
	if btn1.Position().Y != btn2.Position().Y || btn2.Position().Y != btn3.Position().Y {
		t.Errorf("expected all 3 buttons on row 0 at 500px, got y: %v, %v, %v", btn1.Position().Y, btn2.Position().Y, btn3.Position().Y)
	}

	// At 350px, only 2 fit in row 0, btn3 must wrap to row 1
	flow.Layout(btns, fyne.NewSize(350, 100))
	if btn1.Position().Y != btn2.Position().Y {
		t.Errorf("expected btn1 and btn2 on row 0, got y1=%v, y2=%v", btn1.Position().Y, btn2.Position().Y)
	}
	if btn3.Position().Y <= btn1.Position().Y {
		t.Errorf("expected btn3 to wrap to row 1 at 350px, got y3=%v, y1=%v", btn3.Position().Y, btn1.Position().Y)
	}
}

