package layout

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// BuildDemoScreen constructs a full-featured demonstration screen showcasing
// the responsive layout engine across Compact, Medium, Expanded, and awkward screen dimensions.
func BuildDemoScreen(win fyne.Window) fyne.CanvasObject {
	engine := DefaultEngine

	// Live Viewport HUD label
	hudLabel := widget.NewLabel("Viewport: Initializing...")
	hudLabel.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}

	var shell *ResponsiveAppShell

	// 1. Header with Viewport Stats & Window Presets
	menuBtn := widget.NewButtonWithIcon("", theme.MenuIcon(), func() {
		if shell != nil {
			shell.ToggleDrawer()
		}
	})

	titleLabel := widget.NewLabel("Campus OS Layout Engine")
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Preset resolution buttons for quick testing
	preset360 := widget.NewButton("Phone (360x740)", func() {
		if win != nil {
			win.Resize(fyne.NewSize(360, 740))
		}
	})
	preset768 := widget.NewButton("Tablet (768x900)", func() {
		if win != nil {
			win.Resize(fyne.NewSize(768, 900))
		}
	})
	preset1280 := widget.NewButton("Laptop (1280x720)", func() {
		if win != nil {
			win.Resize(fyne.NewSize(1280, 720))
		}
	})
	presetAwkward := widget.NewButton("Awkward (583x420)", func() {
		if win != nil {
			win.Resize(fyne.NewSize(583, 420))
		}
	})

	presetRow := Flow(6, preset360, preset768, preset1280, presetAwkward)

	headerTop := Row(12, menuBtn, titleLabel, Spacer(), hudLabel)
	headerContainer := Column(6, headerTop, presetRow)

	// 2. Desktop Sidebar Navigation
	sideNavItems := []fyne.CanvasObject{
		widget.NewLabelWithStyle("NAVIGATION", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewButtonWithIcon("Dashboard", theme.HomeIcon(), func() {}),
		widget.NewButtonWithIcon("Admissions", theme.AccountIcon(), func() {}),
		widget.NewButtonWithIcon("Ledger", theme.DocumentIcon(), func() {}),
		widget.NewButtonWithIcon("Governance", theme.CheckButtonCheckedIcon(), func() {}),
		widget.NewButtonWithIcon("Settings", theme.SettingsIcon(), func() {}),
	}
	sidebarContent := Column(8, sideNavItems...)
	sidebarContainer := container.NewPadded(sidebarContent)

	// 3. Mobile Slide-in Drawer
	drawerItems := []fyne.CanvasObject{
		widget.NewLabelWithStyle("CAMPUS OS MENU", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewButtonWithIcon("Dashboard", theme.HomeIcon(), func() {
			if shell != nil {
				shell.CloseDrawer()
			}
		}),
		widget.NewButtonWithIcon("Admissions", theme.AccountIcon(), func() {
			if shell != nil {
				shell.CloseDrawer()
			}
		}),
		widget.NewButtonWithIcon("Ledger", theme.DocumentIcon(), func() {
			if shell != nil {
				shell.CloseDrawer()
			}
		}),
		widget.NewButtonWithIcon("Close Drawer", theme.CancelIcon(), func() {
			if shell != nil {
				shell.CloseDrawer()
			}
		}),
	}
	drawerBody := container.NewPadded(Column(10, drawerItems...))

	// Semi-transparent backdrop + Drawer panel
	backdrop := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 160})
	drawerPanelBg := canvas.NewRectangle(color.NRGBA{R: 24, G: 32, B: 42, A: 255})
	drawerPanel := container.NewStack(drawerPanelBg, drawerBody)

	// Drawer docked to left (240px wide) with backdrop
	drawerOverlay := container.NewBorder(nil, nil, container.NewGridWrap(fyne.NewSize(240, 1000), drawerPanel), nil, backdrop)

	// 4. Mobile Bottom Navigation Bar (Compact View)
	bottomNav := container.NewGridWithColumns(4,
		widget.NewButtonWithIcon("Home", theme.HomeIcon(), func() {}),
		widget.NewButtonWithIcon("Students", theme.AccountIcon(), func() {}),
		widget.NewButtonWithIcon("Ledger", theme.DocumentIcon(), func() {}),
		widget.NewButtonWithIcon("More", theme.MenuIcon(), func() {
			if shell != nil {
				shell.ToggleDrawer()
			}
		}),
	)

	// 5. Main Workspace Content (Cards, Grid, Form, Flow)
	// Metric Cards Grid
	card1 := widget.NewCard("Scholars", "1,420 Enrolled", widget.NewLabel("98.4% Active Status"))
	card2 := widget.NewCard("Departments", "8 Active", widget.NewLabel("34 Semesters"))
	card3 := widget.NewCard("Approvals", "3 Pending", widget.NewLabel("Chairperson Action"))
	card4 := widget.NewCard("Staff", "142 Educators", widget.NewLabel("Full Attendance"))

	kpiGrid := Grid(GridOptions{
		MinItemWidth:  240,
		Gap:           12,
		UniformHeight: true,
	}, card1, card2, card3, card4)

	// Responsive Form Controls
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Enter full academic name...")

	roleSelect := widget.NewSelect([]string{"Student", "Faculty", "Registrar", "Dean", "Super Admin"}, func(s string) {})
	roleSelect.SetSelected("Student")

	filterCheck := widget.NewCheck("Only active enrollments", func(b bool) {})
	actionBtn := widget.NewButtonWithIcon("Submit Audit Entry", theme.ConfirmIcon(), func() {})

	formCard := widget.NewCard("Admission Record Filter", "Constraint-Based Adaptive Form",
		Column(10,
			Row(8, widget.NewLabel("Scholar Name:"), FlexItem(nameEntry, 1.0)),
			Row(8, widget.NewLabel("Role Category:"), FlexItem(roleSelect, 1.0)),
			Row(8, filterCheck, Spacer(), actionBtn),
		),
	)

	// Flow Tags Layout
	tag1 := widget.NewButton("B.Tech CSE", func() {})
	tag2 := widget.NewButton("Semester 3", func() {})
	tag3 := widget.NewButton("Lateral Entry", func() {})
	tag4 := widget.NewButton("Hostel Block A", func() {})
	tag5 := widget.NewButton("Section B", func() {})
	tag6 := widget.NewButton("Attendance: 92%", func() {})
	tag7 := widget.NewButton("AKTU Roll: 2200320100123", func() {})

	tagsFlow := Flow(8, tag1, tag2, tag3, tag4, tag5, tag6, tag7)
	tagsCard := widget.NewCard("Dossier Attributes (Flow Layout)", "Items wrap to new rows dynamically", tagsFlow)

	// Information Banner
	infoBanner := widget.NewLabel(
		"Notice: Drag this window's border to see live continuous layout adaptation.\n" +
			"• < 600px: Compact mode with bottom navigation.\n" +
			"• 600-900px: Medium mode with header drawer toggle.\n" +
			"• > 900px: Expanded mode with docked left sidebar.\n" +
			"• Content is capped at 1200px max width and auto-centered on ultrawide monitors.",
	)
	infoBanner.Wrapping = fyne.TextWrapWord

	workspaceContent := Column(16,
		infoBanner,
		widget.NewLabelWithStyle("INSTITUTIONAL METRICS", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		kpiGrid,
		formCard,
		tagsCard,
	)

	// Wrap in fluid container capped at MaxWidth 1200 with AutoCenter
	fluidContainer := Container(workspaceContent, ContainerOptions{
		MaxWidth:      1200,
		AutoCenter:    true,
		Padding:       16,
		ResponsivePad: true,
	})

	scrollableWorkspace := container.NewVScroll(fluidContainer)

	// 6. Build App Shell
	shell = engine.AppShell(AppShellOptions{
		Header:       container.NewPadded(headerContainer),
		Sidebar:      sidebarContainer,
		Content:      scrollableWorkspace,
		BottomNav:    bottomNav,
		Drawer:       drawerOverlay,
		SidebarWidth: 220,
	})

	// Live HUD updater on resize
	shell.container.SetOnSizeChanged(func(size fyne.Size) {
		vp := engine.Viewport(size)
		cols := Columns(size.Width, 240, 12)
		hudLabel.SetText(fmt.Sprintf("%s | %s | %.0fx%.0f (AR: %.2f) | Grid Cols: %d",
			vp.SizeClass(), vp.HeightClass(), size.Width, size.Height, vp.AspectRatio(), cols))
	})

	return shell.Root()
}
