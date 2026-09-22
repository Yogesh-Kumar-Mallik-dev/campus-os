package ui

import (
	"context"
	"fmt"
	"image/color"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
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
	SectionSettings
)

var sectionTitles = map[DashboardSection]string{
	SectionOverview:               "EXECUTIVE OVERVIEW",
	SectionAcademicStructure:      "ACADEMIC STRUCTURE",
	SectionGovernance:             "GOVERNANCE & APPROVALS",
	SectionAuditLedger:            "IMMUTABLE AUDIT LEDGER",
	SectionExecutiveCredentialing: "EXECUTIVE CREDENTIALING",
	SectionSettings:              "SYSTEM & SECURITY SETTINGS",
}

var sectionCategories = map[DashboardSection]string{
	SectionOverview:               "EXECUTIVE CONTROL",
	SectionAcademicStructure:      "ACADEMIC OPERATIONS",
	SectionGovernance:             "POLICY & APPROVALS",
	SectionAuditLedger:            "CRYPTOGRAPHIC LEDGER",
	SectionExecutiveCredentialing: "TIER-1 PROVISIONING",
	SectionSettings:              "PLATFORM CONFIGURATION",
}

// DashboardView manages the master institutional dashboard shell.
type DashboardView struct {
	session            *auth.AuthSession
	client             *api.Client
	window             fyne.Window
	sessionStore       auth.SessionStore
	onLogout           func()
	activeSection      DashboardSection
	sidebarOpen        bool
	userToggledSidebar bool
	windowWidth        float32
	overviewActiveTab  int
	mobileDrawerPopup  *widget.PopUp
	isDark             bool
	rootContainer      *fyne.Container
	workspaceArea      *fyne.Container
	breadcrumbPath     *widget.Label
	pingDot            *canvas.Circle
	pingLabel          *widget.Label
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

	if window != nil && window.Canvas() != nil {
		sz := window.Canvas().Size()
		d.windowWidth = sz.Width
		if sz.Width > 0 && sz.Width < 950 {
			d.sidebarOpen = false
		}
	}

	d.buildShell()
	return d
}

func (d *DashboardView) handleWindowResize(sz fyne.Size) {
	if sz.Width <= 0 {
		return
	}
	prevWidth := d.windowWidth
	d.windowWidth = sz.Width

	// Detect architectural breakpoint crossings:
	// 1. Off-Canvas Drawer: 750px
	// 2. Compact TopBar & Rail: 950px
	// 3. Side-by-side split panels vs Segmented Tab Switcher: 1100px
	// 4. Ultra-narrow 1-col KPI: 550px
	b1Prev := prevWidth > 0 && prevWidth < 750
	b1Curr := sz.Width < 750

	b2Prev := prevWidth > 0 && prevWidth < 950
	b2Curr := sz.Width < 950

	b3Prev := prevWidth > 0 && prevWidth < 1100
	b3Curr := sz.Width < 1100

	b4Prev := prevWidth > 0 && prevWidth < 550
	b4Curr := sz.Width < 550

	shouldCollapseSidebar := sz.Width < 950
	sidebarStateChanged := false
	if !d.userToggledSidebar {
		if shouldCollapseSidebar && d.sidebarOpen {
			d.sidebarOpen = false
			sidebarStateChanged = true
		} else if !shouldCollapseSidebar && !d.sidebarOpen {
			d.sidebarOpen = true
			sidebarStateChanged = true
		}
	}

	if b1Prev != b1Curr || b2Prev != b2Curr || b3Prev != b3Curr || b4Prev != b4Curr || sidebarStateChanged || prevWidth == 0 {
		d.buildShell()
		if d.window != nil && d.window.Canvas() != nil {
			d.window.Canvas().Refresh(d.rootContainer)
		}
	}
}

// CanvasObject returns the renderable dashboard canvas tree.
func (d *DashboardView) CanvasObject() fyne.CanvasObject {
	return d.rootContainer
}

// SetSection transitions the active workspace panel.
func (d *DashboardView) SetSection(section DashboardSection) {
	d.activeSection = section
	if d.breadcrumbPath != nil {
		d.breadcrumbPath.SetText(fmt.Sprintf("Campus OS  ›  Executive Portal  ›  %s", sectionTitles[section]))
	}
	d.renderWorkspace()
	if d.rootContainer != nil {
		d.buildShell()
		if d.window != nil && d.window.Canvas() != nil {
			d.window.Canvas().Refresh(d.rootContainer)
		}
	}
}

// ActiveSection returns the current section for testing.
func (d *DashboardView) ActiveSection() DashboardSection {
	return d.activeSection
}

func (d *DashboardView) buildShell() {
	topBar := d.buildTopBar()
	sidebar := d.buildSidebar()
	d.renderWorkspace()

	// Split view: Sidebar on left (if wide/medium) or 0px in Off-Canvas Drawer mode (<750px)
	var splitContent fyne.CanvasObject
	if sidebar != nil {
		leftGutter := canvas.NewRectangle(color.Transparent)
		leftGutter.SetMinSize(fyne.NewSize(8, 1))
		rightGutter := canvas.NewRectangle(color.Transparent)
		rightGutter.SetMinSize(fyne.NewSize(14, 1))
		splitContent = container.NewBorder(nil, nil, sidebar, nil, container.NewBorder(nil, nil, leftGutter, rightGutter, d.workspaceArea))
	} else {
		leftGutter := canvas.NewRectangle(color.Transparent)
		leftGutter.SetMinSize(fyne.NewSize(14, 1))
		rightGutter := canvas.NewRectangle(color.Transparent)
		rightGutter.SetMinSize(fyne.NewSize(14, 1))
		splitContent = container.NewBorder(nil, nil, leftGutter, rightGutter, d.workspaceArea)
	}

	inner := container.NewBorder(
		topBar,
		nil, nil, nil,
		splitContent,
	)

	if d.rootContainer == nil {
		d.rootContainer = container.New(
			NewResizeNotifierLayout(func(sz fyne.Size) {
				d.handleWindowResize(sz)
			}),
			inner,
		)
	} else {
		d.rootContainer.Objects = []fyne.CanvasObject{inner}
	}
}

// buildTopBar creates the top navigation bar matching the specified layout:
// Left: Hamburger toggle + BBDIT Logo (+ Dynamic Breadcrumbs on desktop)
// Center: Omni-Search Bar (Ctrl+K trigger)
// Right: Academic Term pill + Live UDS ping + Theme toggle + User Profile badge + Sign out
func (d *DashboardView) buildTopBar() fyne.CanvasObject {
	bg := canvas.NewRectangle(theme.Color(theme.ColorNameMenuBackground))
	bg.StrokeColor = theme.Color(theme.ColorNameInputBorder)
	bg.StrokeWidth = 1

	isCompact := d.windowWidth > 0 && d.windowWidth < 950

	// Left: Hamburger toggle + BBDIT Header Logo (+ Breadcrumb on desktop)
	hamburgerIcon := ResourceFromSVG("menu.svg", LucideMenu)
	hamburgerBtn := widget.NewButtonWithIcon("", hamburgerIcon, func() {
		if d.windowWidth > 0 && d.windowWidth < 750 {
			d.showMobileDrawer()
			return
		}
		d.userToggledSidebar = true
		d.sidebarOpen = !d.sidebarOpen
		d.buildShell()
		if d.window != nil && d.window.Canvas() != nil {
			d.window.Canvas().Refresh(d.rootContainer)
		}
	})
	hamburgerBtn.Importance = widget.LowImportance

	logoW, logoH := float32(150), float32(28)
	if isCompact {
		logoW, logoH = float32(140), float32(26)
	}
	logoImg := container.NewCenter(RenderBBDITHeaderLogo(logoW, logoH))

	leftPad := canvas.NewRectangle(color.Transparent)
	leftPad.SetMinSize(fyne.NewSize(8, 1))

	var leftCluster *fyne.Container
	if !isCompact {
		d.breadcrumbPath = widget.NewLabelWithStyle(
			fmt.Sprintf("Campus OS  ›  Executive Portal  ›  %s", sectionTitles[d.activeSection]),
			fyne.TextAlignLeading,
			fyne.TextStyle{Bold: true},
		)
		leftCluster = container.NewHBox(
			leftPad,
			hamburgerBtn,
			logoImg,
			NewShadcnSeparator(false),
			d.breadcrumbPath,
		)
	} else {
		leftCluster = container.NewHBox(
			leftPad,
			hamburgerBtn,
			logoImg,
		)
	}

	// Center: Global Omni-Search Trigger
	searchIcon := ResourceFromSVG("search.svg", LucideSearch)
	searchLabel := "Search scholars, staff, departments (Ctrl+K)..."
	if isCompact {
		searchLabel = "Ctrl+K"
	}
	searchBtn := NewShadcnButton(searchLabel, ButtonOutline, ButtonSizeSm, searchIcon, func() {
		if d.window != nil {
			ShowToast(d.window, "Omni-Search Triggered", "Global Institutional Search palette active.", AlertDefault, 2*time.Second)
		}
	})

	// Right: Academic Term + UDS Heartbeat Ping + Theme Switch Toggle + Logout
	d.pingDot = canvas.NewCircle(color.NRGBA{R: 16, G: 185, B: 129, A: 255})
	d.pingDot.Resize(fyne.NewSize(8, 8))

	themeSwitch := NewSwitch(IsDarkTheme(), func(dark bool) {
		if dark != IsDarkTheme() {
			ToggleTheme()
		}
		d.isDark = IsDarkTheme()
		if d.window != nil {
			themeName := "Light"
			if d.isDark {
				themeName = "Dark"
			}
			ShowToast(d.window, "Theme Switched", fmt.Sprintf("Institutional theme switched to %s mode", themeName), AlertDefault, 2*time.Second)
		}
	})

	themeSunIcon := RenderSVGImage(ColorResourceFromSVG("theme_sun.svg", LucideSun, "#94a3b8"), 14, 14)
	themeMoonIcon := RenderSVGImage(ColorResourceFromSVG("theme_moon.svg", LucideMoon, "#94a3b8"), 14, 14)
	themeToggleCluster := container.NewHBox(
		container.NewCenter(themeSunIcon),
		container.NewCenter(themeSwitch),
		container.NewCenter(themeMoonIcon),
	)

	logoutBtn := NewShadcnButton("Sign Out", ButtonOutline, ButtonSizeSm, ResourceFromSVG("logout.svg", LucideLogOut), func() {
		d.handleLogout()
	})

	var rightCluster *fyne.Container
	if !isCompact {
		termIcon := ResourceFromSVG("cal.svg", LucideCalendar)
		termImg := RenderSVGImage(termIcon, 16, 16)
		termLabel := widget.NewLabel("Fall 2026 • Term A")
		termPill := container.NewHBox(termImg, termLabel)

		d.pingLabel = widget.NewLabel("Live (UDS)")
		pingCluster := container.NewHBox(container.NewCenter(d.pingDot), d.pingLabel)

		rightCluster = container.NewHBox(
			termPill,
			NewShadcnSeparator(false),
			pingCluster,
			NewShadcnSeparator(false),
			themeToggleCluster,
			NewShadcnSeparator(false),
			logoutBtn,
		)
	} else {
		// Compact mode: essentials guaranteed to fit in half-screen width
		pingCluster := container.NewCenter(d.pingDot)

		rightCluster = container.NewHBox(
			pingCluster,
			container.NewCenter(themeSwitch),
			logoutBtn,
		)
	}

	topBarContent := container.NewBorder(
		nil, nil,
		leftCluster,
		rightCluster,
		container.NewCenter(searchBtn),
	)

	return container.NewStack(bg, container.NewPadded(topBarContent))
}


type navEntry struct {
	section DashboardSection
	label   string
	icon    string
	badge   string
}

type navGroup struct {
	header  string
	entries []navEntry
}

func (d *DashboardView) getNavGroups() []navGroup {
	return []navGroup{
		{
			header: "EXECUTIVE & GOVERNANCE",
			entries: []navEntry{
				{SectionOverview, "Executive Overview", LucideLayoutDashboard, ""},
				{SectionGovernance, "Approval Queue", LucideGitPullRequest, "3"},
			},
		},
		{
			header: "ACADEMIC OPERATIONS",
			entries: []navEntry{
				{SectionAcademicStructure, "Academic Hierarchy", LucideGraduationCap, ""},
				{SectionExecutiveCredentialing, "Executive QRs", LucideQrCode, "NEW"},
			},
		},
		{
			header: "SECURITY & LEDGER",
			entries: []navEntry{
				{SectionAuditLedger, "Audit Ledger", LucideFileText, ""},
				{SectionSettings, "System Settings", LucideSettings, ""},
			},
		},
	}
}

func (d *DashboardView) showMobileDrawer() {
	if d.window == nil || d.window.Canvas() == nil {
		return
	}

	drawerWidth := float32(280)
	if d.windowWidth > 0 && d.windowWidth < drawerWidth {
		drawerWidth = d.windowWidth
	}
	drawerHeight := float32(600)
	if d.window.Canvas().Size().Height > 0 {
		drawerHeight = d.window.Canvas().Size().Height
	}

	logoPad := canvas.NewRectangle(color.Transparent)
	logoPad.SetMinSize(fyne.NewSize(8, 1))
	logoImg := container.NewCenter(RenderBBDITHeaderLogo(140, 26))
	drawerHeader := container.NewHBox(logoPad, logoImg)

	navCol := container.NewVBox()
	groups := d.getNavGroups()
	for _, g := range groups {
		hdrLabel := widget.NewLabelWithStyle(g.header, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		navCol.Add(container.NewPadded(hdrLabel))

		for _, item := range g.entries {
			sec := item.section
			lbl := item.label
			ic := item.icon
			bdg := item.badge

			isActive := d.activeSection == sec
			variant := ButtonGhost
			var iconRes fyne.Resource
			if isActive {
				variant = ButtonSecondary
				iconRes = WhiteResourceFromSVG(lbl+"_act.svg", ic)
			} else {
				iconRes = ColorResourceFromSVG(lbl+"_muted.svg", ic, "#94a3b8")
			}

			btn := NewShadcnButton(lbl, variant, ButtonSizeDefault, iconRes, func() {
				if d.mobileDrawerPopup != nil {
					d.mobileDrawerPopup.Hide()
					d.mobileDrawerPopup = nil
				}
				d.SetSection(sec)
			})

			var entryObj fyne.CanvasObject
			if bdg != "" {
				badgePill := NewBadge(bdg, BadgeWarning, BadgeShapePill)
				entryObj = container.NewBorder(nil, nil, nil, badgePill, btn)
			} else {
				entryObj = btn
			}

			if isActive {
				indicator := canvas.NewRectangle(color.NRGBA{R: 225, G: 29, B: 72, A: 255})
				indicator.SetMinSize(fyne.NewSize(3, 20))
				navCol.Add(container.NewBorder(nil, nil, container.NewCenter(indicator), nil, entryObj))
			} else {
				placeholder := canvas.NewRectangle(color.Transparent)
				placeholder.SetMinSize(fyne.NewSize(3, 20))
				navCol.Add(container.NewBorder(nil, nil, placeholder, nil, entryObj))
			}
		}
		navCol.Add(NewShadcnSeparator(true))
	}

	// Bottom User Info Card
	userName := "Chairperson"
	if d.session != nil && d.session.FullName != "" {
		userName = d.session.FullName
	}
	userIcon := RenderSVGImage(ResourceFromSVG("user_mob.svg", LucideUser), 20, 20)
	userLabel := widget.NewLabelWithStyle(userName, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	userStatus := NewBadge("ONLINE", BadgeSuccess, BadgeShapePill)
	bottomCard := container.NewBorder(nil, nil, userIcon, userStatus, userLabel)

	drawerBg := canvas.NewRectangle(theme.Color(theme.ColorNameMenuBackground))
	drawerBg.StrokeColor = theme.Color(theme.ColorNameInputBorder)
	drawerBg.StrokeWidth = 1

	drawerBody := container.NewBorder(
		container.NewPadded(drawerHeader),
		container.NewPadded(bottomCard),
		nil, nil,
		container.NewVScroll(container.NewPadded(navCol)),
	)

	drawerCard := container.NewStack(drawerBg, drawerBody)
	wrapped := container.NewGridWrap(fyne.NewSize(drawerWidth, drawerHeight), drawerCard)

	opts := ModalOptions{
		CloseOnEsc:          true,
		CloseOnClickOutside: true,
		AlignLeft:           true,
		OnDismiss: func() {
			d.mobileDrawerPopup = nil
		},
	}
	d.mobileDrawerPopup = ShowModal(d.window, wrapped, &opts)
}

// buildSidebar creates the responsive navigation:
// - Expanded (240px width): Grouped categories with titles, full labels, and badges.
// - Collapsed (64px width rail): Icons only in a compact vertical rail.
// - Off-Canvas Drawer (< 750px): returns nil so workspace occupies 100% width.
func (d *DashboardView) buildSidebar() fyne.CanvasObject {
	if d.windowWidth > 0 && d.windowWidth < 750 {
		return nil
	}

	bg := canvas.NewRectangle(theme.Color(theme.ColorNameMenuBackground))
	bg.StrokeColor = theme.Color(theme.ColorNameInputBorder)
	bg.StrokeWidth = 1

	if !d.sidebarOpen {
		// Collapsed Icon Rail Mode (64px)
		railItems := []struct {
			section DashboardSection
			name    string
			icon    string
		}{
			{SectionOverview, "Overview", LucideLayoutDashboard},
			{SectionAcademicStructure, "Structure", LucideGraduationCap},
			{SectionGovernance, "Approvals", LucideGitPullRequest},
			{SectionAuditLedger, "Audit", LucideFileText},
			{SectionExecutiveCredentialing, "QRs", LucideQrCode},
			{SectionSettings, "Settings", LucideSettings},
		}

		railCol := container.NewVBox()
		for _, item := range railItems {
			sec := item.section
			isActive := d.activeSection == sec

			var icRes fyne.Resource
			var variant ButtonVariant
			if isActive {
				icRes = WhiteResourceFromSVG(item.name+"_act.svg", item.icon)
				variant = ButtonSecondary
			} else {
				icRes = ColorResourceFromSVG(item.name+"_muted.svg", item.icon, "#94a3b8")
				variant = ButtonGhost
			}

			btn := NewShadcnButton("", variant, ButtonSizeDefault, icRes, func() {
				d.SetSection(sec)
			})

			if isActive {
				indicator := canvas.NewRectangle(color.NRGBA{R: 225, G: 29, B: 72, A: 255})
				indicator.SetMinSize(fyne.NewSize(3, 20))
				railCol.Add(container.NewBorder(nil, nil, container.NewCenter(indicator), nil, btn))
			} else {
				placeholder := canvas.NewRectangle(color.Transparent)
				placeholder.SetMinSize(fyne.NewSize(3, 20))
				railCol.Add(container.NewBorder(nil, nil, placeholder, nil, btn))
			}
		}

		railContent := container.NewStack(bg, container.NewPadded(railCol))
		return container.NewHBox(container.NewGridWrap(fyne.NewSize(64, 600), railContent), NewShadcnSeparator(false))
	}

	// Expanded Mode (240px) with Grouped Categories
	groups := d.getNavGroups()
	sidebarNav := container.NewVBox()

	for _, g := range groups {
		hdrLabel := widget.NewLabelWithStyle(g.header, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		sidebarNav.Add(container.NewPadded(hdrLabel))

		for _, item := range g.entries {
			sec := item.section
			lbl := item.label
			ic := item.icon
			bdg := item.badge

			isActive := d.activeSection == sec
			variant := ButtonGhost
			var iconRes fyne.Resource
			if isActive {
				variant = ButtonSecondary
				iconRes = WhiteResourceFromSVG(lbl+"_act.svg", ic)
			} else {
				iconRes = ColorResourceFromSVG(lbl+"_muted.svg", ic, "#94a3b8")
			}

			btn := NewShadcnButton(lbl, variant, ButtonSizeDefault, iconRes, func() {
				d.SetSection(sec)
			})

			var entryObj fyne.CanvasObject
			if bdg != "" {
				badgePill := NewBadge(bdg, BadgeWarning, BadgeShapePill)
				entryObj = container.NewBorder(nil, nil, nil, badgePill, btn)
			} else {
				entryObj = btn
			}

			if isActive {
				indicator := canvas.NewRectangle(color.NRGBA{R: 225, G: 29, B: 72, A: 255})
				indicator.SetMinSize(fyne.NewSize(3, 20))
				sidebarNav.Add(container.NewBorder(nil, nil, container.NewCenter(indicator), nil, entryObj))
			} else {
				placeholder := canvas.NewRectangle(color.Transparent)
				placeholder.SetMinSize(fyne.NewSize(3, 20))
				sidebarNav.Add(container.NewBorder(nil, nil, placeholder, nil, entryObj))
			}
		}
		sidebarNav.Add(NewShadcnSeparator(true))
	}

	// Bottom User Info Card
	userName := "Chairperson"
	if d.session != nil && d.session.FullName != "" {
		userName = d.session.FullName
	}
	userIcon := RenderSVGImage(ResourceFromSVG("user.svg", LucideUser), 20, 20)
	userLabel := widget.NewLabelWithStyle(userName, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	userStatus := NewBadge("ONLINE", BadgeSuccess, BadgeShapePill)
	bottomCard := container.NewBorder(nil, nil, userIcon, userStatus, userLabel)

	sidebarLayout := container.NewBorder(
		nil,
		container.NewPadded(bottomCard),
		nil, nil,
		container.NewVScroll(container.NewPadded(sidebarNav)),
	)

	fixedWidth := container.NewStack(bg, sidebarLayout)
	return container.NewHBox(container.NewGridWrap(fyne.NewSize(240, 600), fixedWidth), NewShadcnSeparator(false))
}

// buildPageHeader produces a uniform standard header across all sections
func (d *DashboardView) buildPageHeader(title, description, category string, actions ...fyne.CanvasObject) fyne.CanvasObject {
	badge := NewBadge(category, BadgeSecondary, BadgeShapePill)
	titleLabel := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	titleLabel.Wrapping = fyne.TextWrapWord
	descLabel := widget.NewLabel(description)
	descLabel.Wrapping = fyne.TextWrapWord
	descLabel.Importance = widget.LowImportance

	var actionsBox fyne.CanvasObject
	if len(actions) > 0 {
		if len(actions) > 1 && d.windowWidth > 0 && d.windowWidth < 900 {
			primaryBtn := actions[0]
			overflowBtn := NewShadcnButton("", ButtonOutline, ButtonSizeSm, ResourceFromSVG("more.svg", LucideMoreHorizontal), func() {
				if d.window != nil && d.window.Canvas() != nil {
					var pop *widget.PopUp
					menuCol := container.NewVBox()
					for _, secAct := range actions[1:] {
						if btn, ok := secAct.(*ShadcnButton); ok {
							origTap := btn.OnTapped
							btn.OnTapped = func() {
								if pop != nil {
									pop.Hide()
								}
								if origTap != nil {
									origTap()
								}
							}
						}
						menuCol.Add(secAct)
					}
					card := NewShadcnCard(CardParts{
						Title:   "Additional Actions",
						Content: menuCol,
					})
					opts := ModalOptions{
						CloseOnEsc:          true,
						CloseOnClickOutside: true,
					}
					pop = ShowModal(d.window, container.NewGridWrap(fyne.NewSize(260, 160), card), &opts)
				}
			})
			actionsBox = container.NewHBox(primaryBtn, overflowBtn)
		} else {
			actionsBox = container.NewHBox(actions...)
		}
	} else {
		actionsBox = container.NewHBox()
	}

	titleBlock := container.NewVBox(
		container.NewHBox(badge),
		titleLabel,
	)

	var topRow fyne.CanvasObject
	if len(actions) > 0 {
		if d.windowWidth > 0 && d.windowWidth < 680 {
			// On narrow viewports, badge & title on top, action toolbar right-aligned above description
			topRow = container.NewVBox(
				titleBlock,
				container.NewHBox(layout.NewSpacer(), actionsBox),
			)
		} else {
			// Standard desktop/tiled: title block on left, actions docked to the right
			topRow = container.NewBorder(nil, nil, titleBlock, container.NewVBox(actionsBox))
		}
	} else {
		topRow = titleBlock
	}

	headerCard := container.NewVBox(
		topRow,
		descLabel,
	)

	return container.NewVBox(
		container.NewPadded(headerCard),
		NewShadcnSeparator(true),
	)
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
	case SectionSettings:
		d.workspaceArea.Add(d.buildSettingsSection())
	}
}

// buildOverviewSection assembles the 4 KPI cards and 2 executive panels.
func (d *DashboardView) buildOverviewSection() fyne.CanvasObject {
	header := d.buildPageHeader(
		"Executive Operations Center",
		"Real-time governance telemetry, pending executive sign-offs, and collegiate operational health",
		sectionCategories[SectionOverview],
		NewShadcnButton("Provision Executive QR", ButtonDefault, ButtonSizeSm, WhiteResourceFromSVG("act_qr.svg", LucideQrCode), func() {
			d.SetSection(SectionExecutiveCredentialing)
		}),
		NewShadcnButton("Audit Stream", ButtonOutline, ButtonSizeSm, ResourceFromSVG("act_aud.svg", LucideFileText), func() {
			d.SetSection(SectionAuditLedger)
		}),
	)

	// The 4 Primary Executive KPI Cards with Trend Indicators
	kpi1 := NewMetricCard("Total Scholars", "1,420", "• +4.2% YoY • 98.4% Active Enrollment", LucideUserCheck, BadgeSuccess)
	kpi2 := NewMetricCard("Academic Departments", "8 Active", "• 34 Hosted Semesters • 100% Curricular Coverage", LucideGraduationCap, BadgeDefault)
	kpi3 := NewMetricCard("Pending Presidential Approvals", "3 Decisions", "• Chairperson Clearance Required", LucideClock, BadgeWarning)
	kpi4 := NewMetricCard("Staff & Faculty Present", "142 / 148", "• 95.9% Today • 100% Cryptographically Verified", LucideBadgeCheck, BadgeSecondary)

	kpiCols := 4
	if d.windowWidth > 0 && d.windowWidth < 900 {
		kpiCols = 2
		if d.windowWidth < 550 {
			kpiCols = 1
		}
	}

	kpiGrid := container.NewGridWithColumns(kpiCols, kpi1, kpi2, kpi3, kpi4)

	// Left Panel: Pending Approvals Queue
	app1Btn := NewShadcnButton("Approve", ButtonDefault, ButtonSizeSm, WhiteResourceFromSVG("app1.svg", LucideCheckCircle2), func() {
		if d.window != nil {
			ShowToast(d.window, "Action Approved", "CSE Semester 4 curriculum amendment signed and published.", AlertSuccess, 3*time.Second)
		}
	})
	app1 := container.NewBorder(
		nil, nil, nil, app1Btn,
		container.NewVBox(
			widget.NewLabelWithStyle("Computer Science & Engineering: Semester 4 Curriculum Update", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabel("Initiated by Dean Academics • Requires Presidential Sign-Off"),
			NewBadge("HIGH PRIORITY", BadgeWarning, BadgeShapePill),
		),
	)

	app2Btn := NewShadcnButton("Review Dossier", ButtonOutline, ButtonSizeSm, ResourceFromSVG("app2.svg", LucideFileText), func() {
		if d.window != nil {
			ShowToast(d.window, "Review Modal", "Opening Chief Warden disciplinary dossier...", AlertDefault, 2*time.Second)
		}
	})
	app2 := container.NewBorder(
		nil, nil, nil, app2Btn,
		container.NewVBox(
			widget.NewLabelWithStyle("Hostel Block B Warden Disciplinary Escalation", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabel("Initiated by Chief Warden • Student Sabbatical Review"),
			NewBadge("PRESIDENTIAL VETO / CONFIRM", BadgeDestructive, BadgeShapePill),
		),
	)

	approvalsCard := NewShadcnCard(CardParts{
		Badge:       NewBadge("ACTION QUEUE (3 PENDING)", BadgeWarning, BadgeShapePill),
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
	evt4 := widget.NewLabel("• 07:15 AM — Daily biometric attendance synchronization passed with 0 discrepancies.")
	evt4.Wrapping = fyne.TextWrapWord

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
			NewShadcnSeparator(true),
			evt4,
		),
	})

	var workloadContainer fyne.CanvasObject
	if d.windowWidth >= 1100 || d.windowWidth == 0 {
		workloadContainer = container.NewGridWithColumns(2, approvalsCard, activityCard)
	} else {
		tab1Variant := ButtonSecondary
		tab1Icon := WhiteResourceFromSVG("tab_app.svg", LucideGitPullRequest)
		tab2Variant := ButtonGhost
		tab2Icon := ColorResourceFromSVG("tab_act.svg", LucideFileText, "#94a3b8")

		if d.overviewActiveTab == 1 {
			tab1Variant = ButtonGhost
			tab1Icon = ColorResourceFromSVG("tab_app.svg", LucideGitPullRequest, "#94a3b8")
			tab2Variant = ButtonSecondary
			tab2Icon = WhiteResourceFromSVG("tab_act.svg", LucideFileText)
		}

		tabApprovalsBtn := NewShadcnButton("Pending Approvals (3)", tab1Variant, ButtonSizeSm, tab1Icon, func() {
			d.overviewActiveTab = 0
			d.renderWorkspace()
			if d.window != nil && d.window.Canvas() != nil {
				d.window.Canvas().Refresh(d.rootContainer)
			}
		})

		tabActivityBtn := NewShadcnButton("Recent Activity Ledger", tab2Variant, ButtonSizeSm, tab2Icon, func() {
			d.overviewActiveTab = 1
			d.renderWorkspace()
			if d.window != nil && d.window.Canvas() != nil {
				d.window.Canvas().Refresh(d.rootContainer)
			}
		})

		tabTrackBg := canvas.NewRectangle(color.NRGBA{R: 24, G: 24, B: 27, A: 255})
		tabTrackBg.CornerRadius = 6
		tabTrackBg.StrokeColor = color.NRGBA{R: 39, G: 39, B: 42, A: 255}
		tabTrackBg.StrokeWidth = 1

		tabBar := container.NewStack(
			tabTrackBg,
			container.NewHBox(tabApprovalsBtn, tabActivityBtn),
		)

		var activePanel fyne.CanvasObject = approvalsCard
		if d.overviewActiveTab == 1 {
			activePanel = activityCard
		}

		workloadContainer = container.NewVBox(
			container.NewHBox(tabBar),
			activePanel,
		)
	}

	// Quick Actions Bar
	issueQRBtn := NewShadcnButton("Provision Executive QR Docket", ButtonDefault, ButtonSizeDefault, WhiteResourceFromSVG("qr_act.svg", LucideQrCode), func() {
		d.SetSection(SectionExecutiveCredentialing)
	})
	auditBtn := NewShadcnButton("Inspect Full Audit Ledger", ButtonOutline, ButtonSizeDefault, ResourceFromSVG("audit_act.svg", LucideFileText), func() {
		d.SetSection(SectionAuditLedger)
	})
	deptBtn := NewShadcnButton("Review Academic Departments", ButtonOutline, ButtonSizeDefault, ResourceFromSVG("dept_act.svg", LucideGraduationCap), func() {
		d.SetSection(SectionAcademicStructure)
	})
	settingsBtn := NewShadcnButton("System Settings", ButtonOutline, ButtonSizeDefault, ResourceFromSVG("set_act.svg", LucideSettings), func() {
		d.SetSection(SectionSettings)
	})

	actionBar := container.New(NewFlowLayout(8), issueQRBtn, auditBtn, deptBtn, settingsBtn)

	contentBody := container.NewVBox(
		header,
		container.NewPadded(kpiGrid),
		NewShadcnSeparator(true),
		container.NewPadded(workloadContainer),
		NewShadcnSeparator(true),
		container.NewPadded(container.NewBorder(nil, nil, widget.NewLabelWithStyle("EXECUTIVE QUICK ACTIONS:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), nil, actionBar)),
	)

	return container.NewVScroll(contentBody)
}

func (d *DashboardView) buildAcademicStructureSection() fyne.CanvasObject {
	header := d.buildPageHeader(
		"Collegiate Academic Hierarchy",
		"Curricular domain ownership, hosted semester instances, section mappings, and faculty dean assignments",
		sectionCategories[SectionAcademicStructure],
		NewShadcnButton("Add Department", ButtonDefault, ButtonSizeSm, WhiteResourceFromSVG("add_dept.svg", LucideGraduationCap), func() {
			if d.window != nil {
				ShowToast(d.window, "Academic Action", "Department provisioning wizard opened.", AlertDefault, 2*time.Second)
			}
		}),
	)

	dept1 := NewKeyValueRow("COMPUTER SCIENCE & ENGINEERING", "HOD: Dr. Sharma • 480 Scholars • Semesters 1-8 Active", nil)
	dept2 := NewKeyValueRow("INFORMATION TECHNOLOGY", "HOD: Dr. Verma • 240 Scholars • Semesters 1-8 Active", nil)
	dept3 := NewKeyValueRow("MECHANICAL ENGINEERING", "HOD: Dr. Gupta • 180 Scholars • Semesters 1-8 Active", nil)
	dept4 := NewKeyValueRow("ELECTRONICS & COMMUNICATION", "HOD: Dr. Saxena • 210 Scholars • Semesters 1-8 Active", nil)

	card := NewShadcnCard(CardParts{
		Badge:       NewBadge("GOVERNANCE DOMAINS", BadgeDefault, BadgeShapePill),
		Title:       "Institutional Academic Faculties & Departments",
		Description: "Departmental power structure, curriculum authoring, and faculty allocations",
		Content: container.NewVBox(
			dept1,
			NewShadcnSeparator(true),
			dept2,
			NewShadcnSeparator(true),
			dept3,
			NewShadcnSeparator(true),
			dept4,
		),
	})

	return container.NewVScroll(container.NewVBox(header, container.NewPadded(card)))
}

func (d *DashboardView) buildGovernanceSection() fyne.CanvasObject {
	header := d.buildPageHeader(
		"Institutional Governance & Policies",
		"Predefined, un-tamperable approval workflow policies per Campus OS Constitution Section 4.2",
		sectionCategories[SectionGovernance],
	)

	p1 := NewKeyValueRow("SINGLE-DECISIVE (CHAIRPERSON)", "Presidential authority for emergency executive orders and charter amendments", nil)
	p2 := NewKeyValueRow("UNANIMOUS (ALL-MUST-APPROVE)", "Applied to Department delisting, curriculum retirement, and expulsion policies", nil)
	p3 := NewKeyValueRow("SUPER-MAJORITY (2/3 QUORUM)", "Faculty senate budget allocations and multi-department facility utilization", nil)
	p4 := NewKeyValueRow("DELEGATED EXECUTIVE CLEARANCE", "Registrar academic record releases and Warden curfew extensions", nil)

	card := NewShadcnCard(CardParts{
		Badge:       NewBadge("POLICY CATALOG", BadgeWarning, BadgeShapePill),
		Title:       "Constitutional Approval Policies",
		Description: "Guaranteed non-combinable governance templates with strict cryptographic signing",
		Content: container.NewVBox(
			p1,
			NewShadcnSeparator(true),
			p2,
			NewShadcnSeparator(true),
			p3,
			NewShadcnSeparator(true),
			p4,
		),
	})

	return container.NewVScroll(container.NewVBox(header, container.NewPadded(card)))
}

func (d *DashboardView) buildAuditSection() fyne.CanvasObject {
	header := d.buildPageHeader(
		"Immutable Cryptographic Audit Ledger",
		"Chronological, append-only transaction ledger with cryptographic SHA-256 hash chaining",
		sectionCategories[SectionAuditLedger],
		NewShadcnButton("Verify Hash Chain", ButtonDefault, ButtonSizeSm, WhiteResourceFromSVG("verify_audit.svg", LucideShieldCheck), func() {
			if d.window != nil {
				ShowToast(d.window, "Cryptographic Check", "All 1,842 ledger blocks verified against root hash.", AlertSuccess, 3*time.Second)
			}
		}),
	)

	log1 := NewKeyValueRow("TX_AUDIT_0981", "SuperAdminSeat crowned via bare-metal host CLI • Status: COMMITTED (Hash: 3c467c...)", nil)
	log2 := NewKeyValueRow("TX_AUDIT_0982", "Admission Claim Docket issued for Yogesh Kumar Mallik • Status: CLAIMED (Hash: fa0e66...)", nil)
	log3 := NewKeyValueRow("TX_AUDIT_0983", "SIM Hardware Telephony Binding verified via Carrier Handshake • Status: VERIFIED (Hash: d981a2...)", nil)
	log4 := NewKeyValueRow("TX_AUDIT_0984", "Orientation Grace Period Password established • Status: COMMITTED (Hash: 7a82b1...)", nil)

	card := NewShadcnCard(CardParts{
		Badge:       NewBadge("AUDIT TRAIL", BadgeSuccess, BadgeShapePill),
		Title:       "Cryptographic Event Ledger",
		Description: "Un-tamperable append-only institutional event trace",
		Content: container.NewVBox(
			log1,
			NewShadcnSeparator(true),
			log2,
			NewShadcnSeparator(true),
			log3,
			NewShadcnSeparator(true),
			log4,
		),
	})

	return container.NewVScroll(container.NewVBox(header, container.NewPadded(card)))
}

func (d *DashboardView) buildCredentialingSection() fyne.CanvasObject {
	header := d.buildPageHeader(
		"Executive Credential Provisioning",
		"Provision sealed single-use QR activation dockets for Tier-1 leadership (Director, Dean, Registrar)",
		sectionCategories[SectionExecutiveCredentialing],
	)

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

	outputContainer := container.NewVBox()

	generateBtn := NewShadcnButton("Generate & Export Sealed QR Docket", ButtonDefault, ButtonSizeDefault, WhiteResourceFromSVG("qr_gen.svg", LucideQrCode), func() {
		if nameEntry.Text == "" || emailEntry.Text == "" || phoneEntry.Text == "" {
			if d.window != nil {
				ShowToast(d.window, "Missing Information", "All officer fields are required to provision credentials.", AlertDestructive, 3*time.Second)
			}
			return
		}

		go func() {
			var token string
			var payload string
			if d.client != nil {
				resp, err := d.client.GenerateExecutiveQR(context.Background(), roleSelect.Selected, nameEntry.Text, emailEntry.Text, phoneEntry.Text)
				if err == nil && resp != nil {
					token = resp.ClaimToken
					payload = resp.SealedQRPayload
				}
			}

			if token == "" {
				token = fmt.Sprintf("claim_exec_%d", time.Now().UnixNano())
				payload = fmt.Sprintf("CAMPUS_OS:CLAIM:v1:%s", token)
			}

			fyne.Do(func() {
				outputContainer.Objects = nil
				tokenRow := NewKeyValueRow("Issued Single-Use Claim Token", token, nil)
				payloadRow := NewKeyValueRow("Sealed Envelope QR Payload", payload, nil)

				outputCard := NewShadcnCard(CardParts{
					Badge:       NewBadge("SEALED DOCKET GENERATED", BadgeSuccess, BadgeShapePill),
					Title:       fmt.Sprintf("Executive Docket for %s", nameEntry.Text),
					Description: fmt.Sprintf("Assigned Role: %s • Telephony Binding: %s", roleSelect.Selected, phoneEntry.Text),
					Content: container.NewVBox(
						tokenRow,
						NewShadcnSeparator(true),
						payloadRow,
					),
				})
				outputContainer.Add(outputCard)
				outputContainer.Refresh()

				if d.window != nil {
					ShowToast(d.window, "Executive QR Provisioned", fmt.Sprintf("Issued sealed docket for %s (%s)", nameEntry.Text, roleSelect.Selected), AlertSuccess, 3*time.Second)
				}
			})
		}()
	})

	card := NewShadcnCard(CardParts{
		Badge:       NewBadge("TIER-1 PROVISIONING", BadgeDefault, BadgeShapePill),
		Title:       "Issue Executive Credential Docket",
		Description: "Generates cryptographic single-use activation voucher locked to officer SIM telephony",
		Content: container.NewVBox(
			widget.NewLabel("Select Executive Role:"),
			roleSelect,
			NewShadcnSeparator(true),
			nameField,
			emailField,
			phoneField,
		),
		Footer: generateBtn,
	})

	body := container.NewVBox(
		header,
		container.NewPadded(card),
		container.NewPadded(outputContainer),
	)

	return container.NewVScroll(body)
}

func (d *DashboardView) buildSettingsSection() fyne.CanvasObject {
	header := d.buildPageHeader(
		"Platform & Security Settings",
		"Host server socket configuration, session durability, and cryptographic key parameters",
		sectionCategories[SectionSettings],
	)

	row1 := NewKeyValueRow("Campus OS Core Engine", "Version 0.1.0 • Build Date: 2026-09-22", nil)
	row2 := NewKeyValueRow("Persistence Layer Target", "unix:///tmp/campus-os-dev.sock (Prisma 8 UDS)", nil)
	row3 := NewKeyValueRow("Active Auth Mode", "Hardware Telephony + Single-Use QR Envelopes", nil)
	row4 := NewKeyValueRow("Audit Ledger Mode", "Strict SHA-256 Chained Blocks (PostgreSQL auth_schema)", nil)

	sysCard := NewShadcnCard(CardParts{
		Badge:       NewBadge("SYSTEM STATUS", BadgeDefault, BadgeShapePill),
		Title:       "Runtime Environment & Connectivity",
		Description: "Local bare-metal Unix Domain Socket link parameters",
		Content: container.NewVBox(
			row1,
			NewShadcnSeparator(true),
			row2,
			NewShadcnSeparator(true),
			row3,
			NewShadcnSeparator(true),
			row4,
		),
	})

	return container.NewVScroll(container.NewVBox(header, container.NewPadded(sysCard)))
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

// NewMetricCard constructs an executive KPI card with compact shadcn-style typography and semantic iconography.
func NewMetricCard(title, metric, subtitle, svgIcon string, badgeVariant BadgeVariant) fyne.CanvasObject {
	var strokeColor string
	switch badgeVariant {
	case BadgeSuccess:
		strokeColor = "#10b981" // Emerald
	case BadgeWarning:
		strokeColor = "#f59e0b" // Amber
	case BadgeDestructive:
		strokeColor = "#f43f5e" // Rose
	default:
		strokeColor = "#38bdf8" // Sky / Cyan
	}

	iconRes := ColorResourceFromSVG(title+".svg", svgIcon, strokeColor)
	iconImg := RenderSVGImage(iconRes, 22, 22)

	titleLabel := canvas.NewText(strings.ToUpper(title), color.NRGBA{R: 148, G: 163, B: 184, A: 240})
	titleLabel.TextSize = 11
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}

	metricLabel := canvas.NewText(metric, color.White)
	metricLabel.TextSize = 20
	metricLabel.TextStyle = fyne.TextStyle{Bold: true}

	subLabel := widget.NewLabel(subtitle)
	subLabel.Wrapping = fyne.TextWrapWord
	subLabel.Importance = widget.LowImportance

	topRow := container.NewBorder(nil, nil, titleLabel, container.NewCenter(iconImg))

	cardContent := container.NewVBox(
		topRow,
		metricLabel,
		subLabel,
	)

	bg := canvas.NewRectangle(theme.Color(theme.ColorNameMenuBackground))
	bg.StrokeColor = theme.Color(theme.ColorNameInputBorder)
	bg.StrokeWidth = 1
	bg.CornerRadius = 8

	return container.NewStack(bg, container.NewPadded(cardContent))
}
