package ui

import (
	"context"
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Yogesh-Kumar-Mallik-dev/campus-os/apps/client/internal/api"
	"github.com/Yogesh-Kumar-Mallik-dev/campus-os/apps/client/internal/auth"
)

// DashboardSection represents an active navigation destination.
type DashboardSection int

const (
	SectionOverview DashboardSection = iota
	SectionAcademicStructure
	SectionGovernance
	SectionAuditLedger
	SectionExecutiveCredentialing
)

var sectionTitles = map[DashboardSection]string{
	SectionOverview:               "EXECUTIVE OVERVIEW",
	SectionAcademicStructure:      "ACADEMIC STRUCTURE",
	SectionGovernance:             "GOVERNANCE & APPROVALS",
	SectionAuditLedger:            "IMMUTABLE AUDIT LEDGER",
	SectionExecutiveCredentialing: "EXECUTIVE CREDENTIALING",
}

// DashboardView manages the master institutional dashboard shell.
type DashboardView struct {
	session        *auth.AuthSession
	client         *api.Client
	window         fyne.Window
	sessionStore   auth.SessionStore
	onLogout       func()
	activeSection  DashboardSection
	sidebarOpen    bool
	isDark         bool
	rootContainer  *fyne.Container
	workspaceArea  *fyne.Container
	titleLabel     *widget.Label
	pingDot        *canvas.Circle
	pingLabel      *widget.Label
}

// BLOCK_UI_DASHBOARD_NEW_001
// Purpose: Constructs a modern institutional DashboardView shell following Design System 2026.
func NewDashboardView(
	session *auth.AuthSession,
	client *api.Client,
	window fyne.Window,
	sessionStore auth.SessionStore,
	onLogout func(),
) *DashboardView {
	d := &DashboardView{
		session:       session,
		client:        client,
		window:        window,
		sessionStore:  sessionStore,
		onLogout:      onLogout,
		activeSection: SectionOverview,
		sidebarOpen:   true,
		isDark:        true,
		workspaceArea: container.NewStack(),
	}
	d.buildShell()
	return d
}

// CanvasObject returns the renderable dashboard canvas tree.
func (d *DashboardView) CanvasObject() fyne.CanvasObject {
	return d.rootContainer
}

// SetSection transitions the active workspace panel.
func (d *DashboardView) SetSection(section DashboardSection) {
	d.activeSection = section
	if d.titleLabel != nil {
		d.titleLabel.SetText(sectionTitles[section])
	}
	d.renderWorkspace()
}

// ActiveSection returns the current section for testing.
func (d *DashboardView) ActiveSection() DashboardSection {
	return d.activeSection
}

func (d *DashboardView) buildShell() {
	topBar := d.buildTopBar()
	sidebar := d.buildSidebar()
	d.renderWorkspace()

	// Split view: Sidebar on left (if open), workspace taking the rest
	splitContent := container.NewBorder(nil, nil, sidebar, nil, d.workspaceArea)

	d.rootContainer = container.NewBorder(
		topBar,
		nil, nil, nil,
		splitContent,
	)
}

// buildTopBar creates the top navigation bar matching the specified layout:
// Left: Hamburger menu + BBDIT Logo
// Center: Active Page Title
// Right: Theme toggle + Live backend ping + Role badge + Log out
func (d *DashboardView) buildTopBar() fyne.CanvasObject {
	bg := canvas.NewRectangle(theme.Color(theme.ColorNameMenuBackground))
	bg.StrokeColor = theme.Color(theme.ColorNameInputBorder)
	bg.StrokeWidth = 1

	// Left: Hamburger + BBDIT Header Logo
	hamburgerIcon := ResourceFromSVG("menu.svg", LucideMenu)
	hamburgerBtn := widget.NewButtonWithIcon("", hamburgerIcon, func() {
		d.sidebarOpen = !d.sidebarOpen
		d.buildShell()
		if d.window != nil && d.window.Canvas() != nil {
			d.window.Canvas().Refresh(d.rootContainer)
		}
	})
	hamburgerBtn.Importance = widget.LowImportance

	logoImg := RenderBBDITHeaderLogo(120, 24)
	leftCluster := container.NewHBox(hamburgerBtn, logoImg)

	// Center: Page Title
	d.titleLabel = widget.NewLabelWithStyle(sectionTitles[d.activeSection], fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// Right: Theme toggle + Backend Ping + Profile Badge + Logout
	themeIcon := ResourceFromSVG("sun.svg", LucideSun)
	themeBtn := widget.NewButtonWithIcon("", themeIcon, func() {
		d.isDark = !d.isDark
		if d.window != nil {
			ShowToast(d.window, "Theme Switched", "Toggled between dark and light institutional theme", AlertDefault, 2*time.Second)
		}
	})
	themeBtn.Importance = widget.LowImportance

	d.pingDot = canvas.NewCircle(color.NRGBA{R: 16, G: 185, B: 129, A: 255})
	d.pingDot.Resize(fyne.NewSize(8, 8))
	d.pingLabel = widget.NewLabel("Live")
	pingCluster := container.NewHBox(container.NewCenter(d.pingDot), d.pingLabel)

	userName := "User"
	roleName := "STUDENT"
	if d.session != nil {
		if d.session.FullName != "" {
			userName = d.session.FullName
		} else if d.session.Username != "" {
			userName = d.session.Username
		}
		if d.session.RoleCode != "" {
			roleName = d.session.RoleCode
		}
	}
	roleBadge := NewBadge(fmt.Sprintf("%s (%s)", userName, roleName), BadgeDefault, BadgeShapePill)

	logoutBtn := NewShadcnButton("Sign Out", ButtonOutline, ButtonSizeSm, ResourceFromSVG("logout.svg", LucideLogOut), func() {
		d.handleLogout()
	})

	rightCluster := container.NewHBox(
		themeBtn,
		NewShadcnSeparator(false),
		pingCluster,
		NewShadcnSeparator(false),
		roleBadge,
		logoutBtn,
	)

	topBarContent := container.NewBorder(
		nil, nil,
		leftCluster,
		rightCluster,
		container.NewCenter(d.titleLabel),
	)

	return container.NewStack(bg, container.NewPadded(topBarContent))
}

// buildSidebar creates the fixed 220px sidebar for desktop navigation.
func (d *DashboardView) buildSidebar() fyne.CanvasObject {
	if !d.sidebarOpen {
		return container.NewHBox() // Collapsed
	}

	bg := canvas.NewRectangle(theme.Color(theme.ColorNameMenuBackground))
	bg.StrokeColor = theme.Color(theme.ColorNameInputBorder)
	bg.StrokeWidth = 1

	navItems := []struct {
		section DashboardSection
		label   string
		icon    string
	}{
		{SectionOverview, "Executive Overview", LucideLayoutDashboard},
		{SectionAcademicStructure, "Academic Structure", LucideGraduationCap},
		{SectionGovernance, "Governance & Approvals", LucideGitPullRequest},
		{SectionAuditLedger, "Audit Ledger", LucideShieldCheck},
		{SectionExecutiveCredentialing, "Executive QRs", LucideQrCode},
	}

	navButtons := container.NewVBox()

	for _, item := range navItems {
		sec := item.section
		lbl := item.label
		ic := item.icon

		btnVariant := ButtonGhost
		if d.activeSection == sec {
			btnVariant = ButtonSecondary
		}

		iconRes := ResourceFromSVG(lbl+".svg", ic)
		btn := NewShadcnButton(lbl, btnVariant, ButtonSizeDefault, iconRes, func() {
			d.SetSection(sec)
		})
		navButtons.Add(btn)
	}

	sidebarContent := container.NewBorder(
		container.NewPadded(widget.NewLabelWithStyle("INSTITUTIONAL NAVIGATION", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})),
		nil, nil, nil,
		container.NewPadded(navButtons),
	)

	fixedWidth := container.NewStack(bg, sidebarContent)
	return container.NewHBox(container.NewGridWrap(fyne.NewSize(220, 600), fixedWidth), NewShadcnSeparator(false))
}

func (d *DashboardView) renderWorkspace() {
	d.workspaceArea.Objects = nil

	switch d.activeSection {
	case SectionOverview:
		d.workspaceArea.Add(d.buildOverviewSection())
	case SectionAcademicStructure:
		d.workspaceArea.Add(d.buildAcademicStructureSection())
	case SectionGovernance:
		d.workspaceArea.Add(d.buildGovernanceSection())
	case SectionAuditLedger:
		d.workspaceArea.Add(d.buildAuditSection())
	case SectionExecutiveCredentialing:
		d.workspaceArea.Add(d.buildCredentialingSection())
	}
}

// buildOverviewSection assembles the 4 KPI cards and 2 executive panels.
func (d *DashboardView) buildOverviewSection() fyne.CanvasObject {
	// The 4 Primary Executive KPI Cards
	kpi1 := NewMetricCard("Total Scholars", "1,420", "98.4% Active Enrollment", LucideUserCheck, BadgeSuccess)
	kpi2 := NewMetricCard("Academic Departments", "8 Active", "34 Hosted Semesters", LucideGraduationCap, BadgeDefault)
	kpi3 := NewMetricCard("Pending Presidential Approvals", "3 Decisions", "Immediate Action Required", LucideClock, BadgeWarning)
	kpi4 := NewMetricCard("Total Staff & Faculty", "142 Active", "100% Cryptographically Verified", LucideBadgeCheck, BadgeSecondary)

	kpiGrid := container.NewGridWithColumns(4, kpi1, kpi2, kpi3, kpi4)

	// Left Panel: Pending Approvals Queue
	app1 := container.NewVBox(
		widget.NewLabelWithStyle("Computer Science & Engineering: Semester 4 Curriculum Update", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Initiated by Dean Academics • Requires Presidential Sign-Off"),
		NewBadge("HIGH PRIORITY", BadgeWarning, BadgeShapePill),
	)
	app2 := container.NewVBox(
		widget.NewLabelWithStyle("Hostel Block B Warden Disciplinary Escalation", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Initiated by Chief Warden • Student Disciplinary Sabbatical Review"),
		NewBadge("PRESIDENTIAL VETO / CONFIRM", BadgeDestructive, BadgeShapePill),
	)
	approvalsCard := NewShadcnCard(CardParts{
		Badge:       NewBadge("ACTION QUEUE", BadgeWarning, BadgeShapePill),
		Title:       "Pending Executive Approvals",
		Description: "Decisive institutional actions requiring single-signature Chairperson clearance",
		Content: container.NewVBox(
			app1,
			NewShadcnSeparator(true),
			app2,
		),
	})

	// Right Panel: Recent Institutional Activity & Audit Stream
	evt1 := widget.NewLabel("• 10:14 AM — Lateral Entry Scholar (Yogesh Kumar Mallik) claimed admission docket via SIM Slot 1")
	evt1.Wrapping = fyne.TextWrapWord
	evt2 := widget.NewLabel("• 09:45 AM — Academic Registrar published AKTU Semester 3 result manifest (382 records verified)")
	evt2.Wrapping = fyne.TextWrapWord
	evt3 := widget.NewLabel("• 08:30 AM — Automated night audit completed. All UDS persistence sockets intact.")
	evt3.Wrapping = fyne.TextWrapWord

	activityCard := NewShadcnCard(CardParts{
		Badge:       NewBadge("AUDIT TRACE", BadgeDefault, BadgeShapePill),
		Title:       "Recent Institutional Activity Stream",
		Description: "Real-time cryptographically chained event ledger",
		Content: container.NewVBox(
			evt1,
			NewShadcnSeparator(true),
			evt2,
			NewShadcnSeparator(true),
			evt3,
		),
	})

	panelsGrid := container.NewGridWithColumns(2, approvalsCard, activityCard)

	// Quick Actions Bar
	issueQRBtn := NewShadcnButton("Provision Executive QR Docket", ButtonDefault, ButtonSizeDefault, WhiteResourceFromSVG("qr_act.svg", LucideQrCode), func() {
		d.SetSection(SectionExecutiveCredentialing)
	})
	auditBtn := NewShadcnButton("Inspect Full Audit Ledger", ButtonOutline, ButtonSizeDefault, ResourceFromSVG("audit_act.svg", LucideShieldCheck), func() {
		d.SetSection(SectionAuditLedger)
	})
	deptBtn := NewShadcnButton("Review Academic Departments", ButtonOutline, ButtonSizeDefault, ResourceFromSVG("dept_act.svg", LucideGraduationCap), func() {
		d.SetSection(SectionAcademicStructure)
	})

	actionBar := container.NewHBox(issueQRBtn, auditBtn, deptBtn)

	return container.NewVBox(
		container.NewPadded(kpiGrid),
		NewShadcnSeparator(true),
		container.NewPadded(panelsGrid),
		NewShadcnSeparator(true),
		container.NewPadded(container.NewBorder(nil, nil, widget.NewLabelWithStyle("EXECUTIVE QUICK ACTIONS:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), nil, actionBar)),
	)
}

func (d *DashboardView) buildAcademicStructureSection() fyne.CanvasObject {
	dept1 := NewKeyValueRow("COMPUTER SCIENCE & ENGINEERING", "HOD: Dr. Sharma • 480 Scholars • Semesters 1-8 Active", nil)
	dept2 := NewKeyValueRow("INFORMATION TECHNOLOGY", "HOD: Dr. Verma • 240 Scholars • Semesters 1-8 Active", nil)
	dept3 := NewKeyValueRow("MECHANICAL ENGINEERING", "HOD: Dr. Gupta • 180 Scholars • Semesters 1-8 Active", nil)

	return container.NewPadded(NewShadcnCard(CardParts{
		Badge:       NewBadge("GOVERNANCE", BadgeDefault, BadgeShapePill),
		Title:       "Institutional Academic Faculties & Departments",
		Description: "Curricular domain ownership, hosted semester instances, and faculty assignments",
		Content: container.NewVBox(
			dept1,
			NewShadcnSeparator(true),
			dept2,
			NewShadcnSeparator(true),
			dept3,
		),
	}))
}

func (d *DashboardView) buildGovernanceSection() fyne.CanvasObject {
	p1 := NewKeyValueRow("SINGLE-DECISIVE (CHAIRPERSON)", "Presidential authority for emergency executive orders and charter amendments", nil)
	p2 := NewKeyValueRow("UNANIMOUS (ALL-MUST-APPROVE)", "Applied to Department delisting, curriculum retirement, and expulsion policies", nil)

	return container.NewPadded(NewShadcnCard(CardParts{
		Badge:       NewBadge("WORKFLOW POLICIES", BadgeWarning, BadgeShapePill),
		Title:       "Predefined Approval Policy Catalog",
		Description: "Strict non-combinable governance templates per Campus OS Constitution Section 4.2",
		Content: container.NewVBox(
			p1,
			NewShadcnSeparator(true),
			p2,
		),
	}))
}

func (d *DashboardView) buildAuditSection() fyne.CanvasObject {
	log1 := NewKeyValueRow("TX_AUDIT_0981", "SuperAdminSeat initialized via bare-metal host CLI • Status: COMMITTED", nil)
	log2 := NewKeyValueRow("TX_AUDIT_0982", "Admission Claim Docket issued for Yogesh Kumar Mallik • Status: CLAIMED", nil)

	return container.NewPadded(NewShadcnCard(CardParts{
		Badge:       NewBadge("IMMUTABLE TRACE", BadgeSuccess, BadgeShapePill),
		Title:       "Cryptographic Audit Ledger",
		Description: "Un-tamperable append-only institutional event trace",
		Content: container.NewVBox(
			log1,
			NewShadcnSeparator(true),
			log2,
		),
	}))
}

func (d *DashboardView) buildCredentialingSection() fyne.CanvasObject {
	nameField, nameEntry := NewFormField(FormField{
		Label:       "Executive Officer Legal Full Name",
		Placeholder: "e.g. Dr. Rajesh Sharma",
	})
	emailField, emailEntry := NewFormField(FormField{
		Label:       "Institutional Email Address",
		Placeholder: "e.g. director@campus.edu",
	})
	phoneField, phoneEntry := NewFormField(FormField{
		Label:       "Registered Telephony Contact (SIM Binding)",
		Placeholder: "e.g. +91 98765-43210",
	})

	roleSelect := widget.NewSelect([]string{
		"Director of Institution",
		"Executive Director",
		"Dean of Academics",
		"Academic Registrar",
	}, nil)
	roleSelect.SetSelected("Director of Institution")

	generateBtn := NewShadcnButton("Generate & Export Sealed QR Docket", ButtonDefault, ButtonSizeDefault, WhiteResourceFromSVG("qr_gen.svg", LucideQrCode), func() {
		if nameEntry.Text == "" || emailEntry.Text == "" || phoneEntry.Text == "" {
			if d.window != nil {
				ShowToast(d.window, "Missing Information", "All officer fields are required to provision credentials.", AlertDestructive, 3*time.Second)
			}
			return
		}
		if d.window != nil {
			ShowToast(d.window, "Executive QR Provisioned", fmt.Sprintf("Issued sealed docket for %s (%s)", nameEntry.Text, roleSelect.Selected), AlertSuccess, 3*time.Second)
		}
	})

	return container.NewPadded(NewShadcnCard(CardParts{
		Badge:       NewBadge("TIER-1 PROVISIONING", BadgeDefault, BadgeShapePill),
		Title:       "Executive Credential Provisioning",
		Description: "Generate single-use sealed QR activation dockets for institutional leadership",
		Content: container.NewVBox(
			widget.NewLabel("Select Executive Role:"),
			roleSelect,
			NewShadcnSeparator(true),
			nameField,
			emailField,
			phoneField,
		),
		Footer: generateBtn,
	}))
}

// handleLogout executes clean dual-revocation: clears OS Keyring and invalidates server session.
func (d *DashboardView) handleLogout() {
	if d.client != nil && d.session != nil && d.session.RefreshToken != "" {
		go func() {
			_ = d.client.RevokeSession(context.Background(), d.session.RefreshToken)
		}()
	}
	if d.sessionStore != nil {
		_ = d.sessionStore.Clear()
	}
	if d.onLogout != nil {
		d.onLogout()
	}
}

// NewMetricCard constructs an executive KPI card with luminous glass accenting.
func NewMetricCard(title, metric, subtitle, svgIcon string, badgeVariant BadgeVariant) fyne.CanvasObject {
	iconRes := ResourceFromSVG(title+".svg", svgIcon)
	iconImg := RenderSVGImage(iconRes, 28, 28)

	titleLabel := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{})
	titleLabel.Wrapping = fyne.TextWrapWord

	metricLabel := widget.NewLabelWithStyle(metric, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	subLabel := widget.NewLabel(subtitle)
	subLabel.Wrapping = fyne.TextWrapWord

	topRow := container.NewBorder(nil, nil, nil, iconImg, titleLabel)

	cardContent := container.NewVBox(
		topRow,
		metricLabel,
		NewShadcnSeparator(true),
		subLabel,
	)

	bg := canvas.NewRectangle(theme.Color(theme.ColorNameMenuBackground))
	bg.StrokeColor = theme.Color(theme.ColorNameInputBorder)
	bg.StrokeWidth = 1
	bg.CornerRadius = 8

	return container.NewStack(bg, container.NewPadded(cardContent))
}
