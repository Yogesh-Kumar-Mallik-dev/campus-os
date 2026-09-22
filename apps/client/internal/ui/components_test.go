package ui

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func TestCampusTheme_ColorsAndSizes(t *testing.T) {
	th := NewCampusTheme()

	// Light variant checks
	lightBg := th.Color(theme.ColorNameBackground, theme.VariantLight)
	if lightBg == nil {
		t.Fatal("expected non-nil light background")
	}
	lightFg := th.Color(theme.ColorNameForeground, theme.VariantLight)
	if lightFg == nil {
		t.Fatal("expected non-nil light foreground")
	}
	lightPrimary := th.Color(theme.ColorNamePrimary, theme.VariantLight)
	if lightPrimary == nil {
		t.Fatal("expected non-nil light primary")
	}

	// Dark variant checks
	darkBg := th.Color(theme.ColorNameBackground, theme.VariantDark)
	if darkBg == nil {
		t.Fatal("expected non-nil dark background")
	}
	darkFg := th.Color(theme.ColorNameForeground, theme.VariantDark)
	if darkFg == nil {
		t.Fatal("expected non-nil dark foreground")
	}
	darkPrimary := th.Color(theme.ColorNamePrimary, theme.VariantDark)
	if darkPrimary == nil {
		t.Fatal("expected non-nil dark primary")
	}

	// Interactive & Accessibility state colors
	for _, variant := range []fyne.ThemeVariant{theme.VariantLight, theme.VariantDark} {
		for _, colorName := range []fyne.ThemeColorName{
			theme.ColorNameHover,
			theme.ColorNamePressed,
			theme.ColorNameFocus,
			theme.ColorNameDisabled,
			theme.ColorNameDisabledButton,
			theme.ColorNameSelection,
			theme.ColorNameHyperlink,
			theme.ColorNameScrollBar,
		} {
			c := th.Color(colorName, variant)
			if c == nil {
				t.Fatalf("expected non-nil color for %s in variant %d", colorName, variant)
			}
		}
	}

	// Size checks
	cardRadius := th.Size(theme.SizeNameCardRadius)
	if cardRadius != 10.0 {
		t.Errorf("expected card radius 10.0, got %f", cardRadius)
	}

	btnRadius := th.Size(theme.SizeNameButtonRadius)
	if btnRadius != 8.0 {
		t.Errorf("expected button radius 8.0, got %f", btnRadius)
	}

	headingSize := th.Size(theme.SizeNameHeadingText)
	if headingSize != 19.0 {
		t.Errorf("expected heading size 19.0, got %f", headingSize)
	}
}

func TestComponents_Render(t *testing.T) {
	// Status Pills
	pillVariants := []PillVariant{PillSuccess, PillWarning, PillError, PillInfo, PillNeutral}
	for _, v := range pillVariants {
		p := NewStatusPill("TEST", v)
		if p == nil {
			t.Fatalf("expected non-nil pill for variant %d", v)
		}
	}

	// Page Header with and without pill
	hWithPill := NewPageHeader("Title", "Subtitle", NewStatusPill("STATUS", PillSuccess))
	if hWithPill == nil {
		t.Fatal("expected non-nil page header with pill")
	}

	hNoPill := NewPageHeader("Title", "Subtitle", nil)
	if hNoPill == nil {
		t.Fatal("expected non-nil page header without pill")
	}

	// Styled Card
	c1 := NewStyledCard("Card Title", hNoPill)
	if c1 == nil {
		t.Fatal("expected non-nil styled card with title")
	}

	c2 := NewStyledCard("", hNoPill)
	if c2 == nil {
		t.Fatal("expected non-nil styled card without title")
	}

	// Banner Notice with and without SVG
	banner1 := NewBannerNotice("Notice", "Message", PillWarning, ResourceFromSVG("test.svg", SVGShieldCheck))
	if banner1 == nil {
		t.Fatal("expected non-nil banner with SVG")
	}

	banner2 := NewBannerNotice("Notice", "Message", PillInfo, nil)
	if banner2 == nil {
		t.Fatal("expected non-nil banner without SVG")
	}

	// Step Indicator
	stepper := NewStepIndicator(2, []string{"Step 1", "Step 2", "Step 3"})
	if stepper == nil {
		t.Fatal("expected non-nil step indicator")
	}

	// Top Bar
	topBar := NewTopBar("CAMPUS OS", "Account Activation")
	if topBar == nil {
		t.Fatal("expected non-nil top bar")
	}
	topBarNoRole := NewTopBar("CAMPUS OS", "")
	if topBarNoRole == nil {
		t.Fatal("expected non-nil top bar without role")
	}

	// Metric Cards across semantic variants
	for _, v := range []BadgeVariant{BadgeSuccess, BadgeWarning, BadgeDestructive, BadgeSecondary, BadgeDefault} {
		m := NewMetricCard("Active Users", "1,234", "+5% this week", LucideUserCheck, v)
		if m == nil {
			t.Fatalf("expected non-nil metric card for variant %v", v)
		}
	}

	// Approval Queue Item
	appItem := NewApprovalQueueItem("Title", "Subtitle", "HIGH", BadgeWarning, "Approve", LucideCheckCircle2, true, func() {})
	if appItem == nil {
		t.Fatal("expected non-nil approval queue item")
	}

	// Activity Ledger Item
	actItem := NewActivityLedgerItem("12:00 PM", "Headline", "Detail", LucideShieldCheck, true)
	if actItem == nil {
		t.Fatal("expected non-nil activity ledger item")
	}
}

func TestAssets_LogoRender(t *testing.T) {
	if BBDITLogoResource == nil || len(BBDITLogoResource.Content()) == 0 {
		t.Fatal("expected non-empty BBDIT logo resource")
	}

	logoImg := RenderLogoImage(160, 30)
	if logoImg == nil {
		t.Fatal("expected non-nil logo image")
	}

	headerLogo := RenderBBDITHeaderLogo(140, 24)
	if headerLogo == nil {
		t.Fatal("expected non-nil header logo pill")
	}
}

func TestIcons_SVGResources(t *testing.T) {
	svgList := []struct {
		name string
		svg  string
	}{
		{"bbdit_logo", SVGBBDITLogo},
		{"lucide_scan", LucideScan},
		{"lucide_qr", LucideQrCode},
		{"lucide_sim", LucideSim},
		{"lucide_sim_sec", LucideSimSecondary},
		{"lucide_smartphone", LucideSmartphone},
		{"lucide_shield_check", LucideShieldCheck},
		{"lucide_user_check", LucideUserCheck},
		{"lucide_file_check", LucideFileCheck},
		{"lucide_key_round", LucideKeyRound},
		{"lucide_lock", LucideLock},
		{"lucide_clock", LucideClock},
		{"lucide_fingerprint", LucideFingerprint},
		{"lucide_id_card", LucideIdCard},
		{"lucide_badge_check", LucideBadgeCheck},
		{"lucide_alert_triangle", LucideAlertTriangle},
		{"lucide_alert_circle", LucideAlertCircle},
		{"lucide_check_circle2", LucideCheckCircle2},
		{"lucide_info", LucideInfo},
		{"lucide_flag", LucideFlag},
		{"lucide_chevron_right", LucideChevronRight},
		{"lucide_sparkles", LucideSparkles},
		{"lucide_graduation_cap", LucideGraduationCap},
		{"lucide_image", LucideImage},
		{"lucide_upload", LucideUpload},
		{"lucide_menu", LucideMenu},
		{"lucide_dashboard", LucideLayoutDashboard},
		{"lucide_pull_request", LucideGitPullRequest},
		{"lucide_login", LucideLogIn},
		{"lucide_logout", LucideLogOut},
		{"lucide_sun", LucideSun},
		{"lucide_moon", LucideMoon},
	}

	for _, s := range svgList {
		res := ResourceFromSVG(s.name, s.svg)
		if res == nil || len(res.Content()) == 0 {
			t.Fatalf("failed to create SVG resource for %s", s.name)
		}
		img := RenderSVGImage(res, 24, 24)
		if img == nil {
			t.Fatalf("failed to render SVG image for %s", s.name)
		}
	}
}

func TestResponsiveLayout(t *testing.T) {
	rl := NewResponsiveLayout(0, 0)
	if rl == nil {
		t.Fatal("expected non-nil responsive layout")
	}

	btn := widget.NewButton("Action", nil)
	containerObj := container.New(rl, btn)

	// Test layout calculation on narrow viewport (mobile 360px)
	rl.Layout([]fyne.CanvasObject{btn}, fyne.NewSize(360, 640))
	if btn.Size().Width != 360-24 { // 336px
		t.Errorf("expected mobile width 336, got %f", btn.Size().Width)
	}

	// Test layout calculation on wide viewport (desktop 1200px)
	rl.Layout([]fyne.CanvasObject{btn}, fyne.NewSize(1200, 800))
	if btn.Size().Width != 580 { // clamped to maxWidth 580
		t.Errorf("expected desktop width clamped to 580, got %f", btn.Size().Width)
	}

	// Min size test
	minS := rl.MinSize([]fyne.CanvasObject{btn})
	if minS.Width <= 0 || minS.Height <= 0 {
		t.Errorf("expected positive min size, got %v", minS)
	}

	// Responsive card container
	respContainer := NewResponsiveCardContainer(containerObj, 640)
	if respContainer == nil {
		t.Fatal("expected non-nil responsive card container")
	}
}
