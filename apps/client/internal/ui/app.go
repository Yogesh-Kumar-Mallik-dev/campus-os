package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Yogesh-Kumar-Mallik-dev/campus-os/apps/client/internal/api"
)

// Application encapsulates the Fyne native UI runtime.
type Application struct {
	FyneApp   fyne.App
	Window    fyne.Window
	APIClient *api.Client
}

// BLOCK_UI_APP_NEW_001
// Purpose: Constructs the Campus OS native client application window and layout.
func NewApplication(apiClient *api.Client) *Application {
	a := app.NewWithID("internal.campus-os.client")
	w := a.NewWindow("Campus OS — Institutional Client")

	return &Application{
		FyneApp:   a,
		Window:    w,
		APIClient: apiClient,
	}
}

// BLOCK_UI_APP_BUILD_001
// Purpose: Assembles primary tab navigation and responsive layout views.
func (a *Application) BuildLayout() fyne.CanvasObject {
	// 1. Attendance Tab
	attendanceList := widget.NewList(
		func() int { return 4 },
		func() fyne.CanvasObject {
			return widget.NewLabel("Session Record")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			labels := []string{
				"Today (2026-09-21): Period 1 (DAA - CS401) — [PRESENT]",
				"Yesterday (2026-09-20): Period 2 (OS Lab) — [ON_DUTY - Sanctioned]",
				"2026-09-19: Period 3 (Maths IV) — [PRESENT]",
				"2026-09-18: Period 4 (Database Systems) — [ABSENT]",
			}
			o.(*widget.Label).SetText(labels[i])
		},
	)
	attendanceTab := container.NewBorder(
		widget.NewLabelWithStyle("Time-Series Attendance Ledger", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		nil, nil, nil,
		attendanceList,
	)

	// 2. Hostel Outpass Tab
	destEntry := widget.NewEntry()
	destEntry.SetPlaceHolder("Enter destination (e.g. City Center)")
	reasonEntry := widget.NewEntry()
	reasonEntry.SetPlaceHolder("Purpose of visit / emergency reason")

	outpassStatus := widget.NewLabel("Current Status: No Active Outpass")
	submitBtn := widget.NewButtonWithIcon("Request Digital Outpass", theme.MailSendIcon(), func() {
		if destEntry.Text == "" {
			outpassStatus.SetText("Error: Destination is required")
			return
		}
		outpassStatus.SetText("Request Submitted! Pending Warden Approval -> Gate QR Code")
	})

	outpassTab := container.NewVBox(
		widget.NewLabelWithStyle("Campus Life & Gate Outpass", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Destination:"),
		destEntry,
		widget.NewLabel("Reason:"),
		reasonEntry,
		submitBtn,
		widget.NewSeparator(),
		outpassStatus,
	)

	// 3. Academic Profile Tab
	profileTab := container.NewVBox(
		widget.NewLabelWithStyle("Student Academic Profile", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewForm(
			widget.NewFormItem("Roll Number", widget.NewLabel("23CS1042")),
			widget.NewFormItem("PRN", widget.NewLabel("PRN-2023-8891")),
			widget.NewFormItem("Department", widget.NewLabel("Computer Science & Engineering")),
			widget.NewFormItem("Course", widget.NewLabel("B.Tech CSE (4th Semester, Section B)")),
			widget.NewFormItem("Admission Type", widget.NewLabel("LATERAL_ENTRY (Admitted Sem 3)")),
			widget.NewFormItem("Status", widget.NewLabel("ACTIVE")),
		),
	)

	// 4. Onboarding / QR Activation Tab
	wizard := NewOnboardingWizard(a.APIClient, func() {
		// Switch to attendance tab upon onboarding completion
	})
	onboardTab := container.NewScroll(wizard.CanvasObject())

	// Tab Container
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Attendance", theme.MenuIcon(), attendanceTab),
		container.NewTabItemWithIcon("Hostel Outpass", theme.NavigateNextIcon(), outpassTab),
		container.NewTabItemWithIcon("Academic Profile", theme.AccountIcon(), profileTab),
		container.NewTabItemWithIcon("Activate Account", theme.LoginIcon(), onboardTab),
	)
	tabs.SetTabLocation(container.TabLocationLeading)

	return tabs
}

// Run launches the native Fyne window loop.
func (a *Application) Run() {
	a.Window.SetContent(a.BuildLayout())
	a.Window.Resize(fyne.NewSize(850, 550))
	a.Window.CenterOnScreen()
	a.Window.ShowAndRun()
}
