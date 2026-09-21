package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
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
// Purpose: Constructs the Campus OS native client application window, applies custom design tokens, and sets responsive geometry.
func NewApplication(apiClient *api.Client) *Application {
	a := app.NewWithID("internal.campus-os.client")
	a.Settings().SetTheme(NewCampusTheme())

	w := a.NewWindow("Campus OS — Institutional Client")

	return &Application{
		FyneApp:   a,
		Window:    w,
		APIClient: apiClient,
	}
}

// BLOCK_UI_APP_BUILD_001
// Purpose: Assembles modern card layout, KPI widgets, and responsive tab views.
func (a *Application) BuildLayout() fyne.CanvasObject {
	topBar := NewTopBar("CAMPUS OS", "Student Scholar")

	// -------------------------------------------------------------------------
	// 1. Attendance Tab
	// -------------------------------------------------------------------------
	attHeader := NewPageHeader(
		"Time-Series Attendance Ledger",
		"Biometric, RFID, and faculty roll time-series records",
		NewStatusPill("84.2% OVERALL", PillSuccess),
	)

	statAttended := NewStatCard("Attended", "41", "Sessions present", PillSuccess)
	statMissed := NewStatCard("Absent", "7", "Unexcused cuts", PillError)
	statExcused := NewStatCard("On Duty", "2", "Medical/Hackathon", PillInfo)
	statThreshold := NewStatCard("AKTU Eligibility", "SAFE", "Min 75% requirement", PillSuccess)

	statsGrid := container.NewGridWithColumns(4, statAttended, statMissed, statExcused, statThreshold)

	historyData := []struct {
		subject string
		date    string
		time    string
		status  string
		pill    PillVariant
	}{
		{"DAA (KCS-401) — Design & Analysis of Algorithms", "Today", "09:30 AM", "PRESENT", PillSuccess},
		{"Operating Systems Lab (KCS-451)", "Yesterday", "02:00 PM", "ON_DUTY", PillInfo},
		{"Discrete Mathematics (KAS-402)", "19 Sep", "11:15 AM", "PRESENT", PillSuccess},
		{"Database Management Systems (KCS-403)", "18 Sep", "10:30 AM", "ABSENT", PillError},
		{"Universal Human Values (KVE-401)", "17 Sep", "03:15 PM", "PRESENT", PillSuccess},
	}

	historyList := widget.NewList(
		func() int { return len(historyData) },
		func() fyne.CanvasObject {
			subjText := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
			timeText := widget.NewLabel("")
			pillPlaceholder := NewStatusPill("STATUS", PillNeutral)
			return container.NewBorder(nil, nil, nil, pillPlaceholder, container.NewVBox(subjText, timeText))
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			item := historyData[i]
			border := o.(*fyne.Container)
			vbox := border.Objects[0].(*fyne.Container)
			vbox.Objects[0].(*widget.Label).SetText(item.subject)
			vbox.Objects[1].(*widget.Label).SetText(item.date + " • " + item.time)
			border.Objects[1] = NewStatusPill(item.status, item.pill)
		},
	)

	historyCard := NewStyledCard("Recent Attendance Activity", container.NewGridWrap(fyne.NewSize(760, 220), historyList))
	attendanceTab := container.NewVBox(attHeader, statsGrid, historyCard)

	// -------------------------------------------------------------------------
	// 2. Hostel Outpass Tab
	// -------------------------------------------------------------------------
	outpassHeader := NewPageHeader(
		"Campus Life & Digital Outpass",
		"Automated mentor routing, parent consent, and biometric gate verification",
		NewStatusPill("NO ACTIVE PASS", PillNeutral),
	)

	destEntry := widget.NewEntry()
	destEntry.SetPlaceHolder("e.g. Anand Vihar Railway Station / Home")
	reasonEntry := widget.NewMultiLineEntry()
	reasonEntry.SetPlaceHolder("Specify urgent reason / semester break departure")

	outpassStatus := widget.NewLabel("No pass requested for current week.")
	submitBtn := widget.NewButtonWithIcon("Submit Outpass Request", theme.MailSendIcon(), func() {
		if destEntry.Text == "" {
			outpassStatus.SetText("Error: Destination is required")
			return
		}
		outpassStatus.SetText("✓ Request submitted! Parent OTP dispatched -> Pending Warden approval.")
	})
	submitBtn.Importance = widget.HighImportance

	formContent := container.NewVBox(
		widget.NewLabel("Destination:"),
		destEntry,
		widget.NewLabel("Departure Reason:"),
		reasonEntry,
		layout.NewSpacer(),
		submitBtn,
		widget.NewSeparator(),
		outpassStatus,
	)

	outpassCard := NewStyledCard("New Pass Application", formContent)
	rulesCard := NewStyledCard("Outpass Guidelines", widget.NewLabel(
		"• Curfew for weekday return is 08:30 PM.\n"+
			"• Night outpasses require one-time SMS confirmation from registered parent mobile.\n"+
			"• Gate security scans your QR token upon departure and return.",
	))

	outpassGrid := container.NewGridWithColumns(2, outpassCard, rulesCard)
	outpassTab := container.NewVBox(outpassHeader, outpassGrid)

	// -------------------------------------------------------------------------
	// 3. Academic Profile Tab
	// -------------------------------------------------------------------------
	profileHeader := NewPageHeader(
		"Academic & Institutional Profile",
		"Affiliated with Dr. A.P.J. Abdul Kalam Technical University (AKTU)",
		NewStatusPill("LATERAL ENTRY", PillWarning),
	)

	namesCard := NewStyledCard("Scholar Identity", container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("Academic Name", widget.NewLabel("Yogesh")),
			widget.NewFormItem("Legal Full Name", widget.NewLabel("Yogesh Kumar Mallik")),
			widget.NewFormItem("Enrollment / Roll No", widget.NewLabel("2200970139001")),
			widget.NewFormItem("Father Name (Legal)", widget.NewLabel("Father Legal Name")),
			widget.NewFormItem("Mother Name (Legal)", widget.NewLabel("Mother Legal Name")),
		),
	))

	progressionCard := NewStyledCard("Curriculum Progression", container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("Department", widget.NewLabel("Computer Science & Engineering")),
			widget.NewFormItem("Admission Mode", widget.NewLabel("Direct Lateral Entry (2nd Year)")),
			widget.NewFormItem("Entry Semester", widget.NewLabel("Semester 3 (Direct Entry)")),
			widget.NewFormItem("Current Semester", widget.NewLabel("Semester 4")),
			widget.NewFormItem("Institutional Status", widget.NewLabel("ACTIVE (Good Standing)")),
		),
	))

	profileGrid := container.NewGridWithColumns(2, namesCard, progressionCard)
	profileTab := container.NewVBox(profileHeader, profileGrid)

	// -------------------------------------------------------------------------
	// 4. Onboarding / QR Activation Tab
	// -------------------------------------------------------------------------
	wizard := NewOnboardingWizard(a.APIClient, func() {})
	onboardHeader := NewPageHeader(
		"First-Time Account Activation",
		"Activate credentials, pair device carrier SIM, and generate Digital Identity",
		NewStatusPill("STEP 1 OF 5", PillInfo),
	)
	onboardTab := container.NewVBox(onboardHeader, wizard.CanvasObject())

	// -------------------------------------------------------------------------
	// Primary App Tabs Navigation
	// -------------------------------------------------------------------------
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Attendance", theme.MenuIcon(), container.NewScroll(attendanceTab)),
		container.NewTabItemWithIcon("Hostel Outpass", theme.NavigateNextIcon(), container.NewScroll(outpassTab)),
		container.NewTabItemWithIcon("Academic Profile", theme.AccountIcon(), container.NewScroll(profileTab)),
		container.NewTabItemWithIcon("Activate Account", theme.LoginIcon(), container.NewScroll(onboardTab)),
	)
	tabs.SetTabLocation(container.TabLocationLeading)

	return container.NewBorder(
		topBar,
		nil, nil, nil,
		tabs,
	)
}

// Run launches the native Fyne window loop.
func (a *Application) Run() {
	a.Window.SetContent(a.BuildLayout())
	a.Window.Resize(fyne.NewSize(960, 640))
	a.Window.CenterOnScreen()
	a.Window.ShowAndRun()
}
