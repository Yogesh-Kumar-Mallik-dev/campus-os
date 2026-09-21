package ui

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func TestBadge_VariantsAndShapes(t *testing.T) {
	variants := []BadgeVariant{
		BadgeDefault,
		BadgeSecondary,
		BadgeDestructive,
		BadgeOutline,
		BadgeSuccess,
		BadgeWarning,
	}

	shapes := []BadgeShape{
		BadgeShapePill,
		BadgeShapeRounded,
	}

	for _, v := range variants {
		for _, s := range shapes {
			b := NewBadge("TEST BADGE", v, s)
			if b == nil {
				t.Fatalf("expected non-nil badge for variant %d, shape %d", v, s)
			}
		}
	}

	// Test InteractiveBadge with full interactive states
	tapped := false
	ib := NewInteractiveBadge("FILTER CHIP", BadgeSecondary, BadgeShapePill, func() {
		tapped = true
	})
	if ib == nil {
		t.Fatal("expected non-nil interactive badge")
	}

	// Tap
	ib.Tapped(nil)
	if !tapped {
		t.Fatal("expected interactive badge to handle tap")
	}
	tapped = false

	// Hover
	ib.MouseIn(nil)
	if !ib.Hovered {
		t.Fatal("expected badge to be hovered")
	}
	ib.MouseOut()
	if ib.Hovered {
		t.Fatal("expected badge not to be hovered")
	}

	// Press
	ib.MouseDown(nil)
	if !ib.Pressed {
		t.Fatal("expected badge to be pressed")
	}
	ib.MouseUp(nil)
	if ib.Pressed {
		t.Fatal("expected badge not to be pressed")
	}

	// Focus & Keyboard
	ib.FocusGained()
	if !ib.Focused {
		t.Fatal("expected badge to be focused")
	}
	ib.TypedRune(' ')
	if !tapped {
		t.Fatal("expected space key to tap badge")
	}
	tapped = false

	ib.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	if !tapped {
		t.Fatal("expected return key to tap badge")
	}
	tapped = false
	ib.FocusLost()

	// Disabled
	ib.Disable()
	if !ib.Disabled() {
		t.Fatal("expected badge to report disabled")
	}
	ib.Tapped(nil)
	if tapped {
		t.Fatal("expected disabled badge to ignore tap")
	}
	ib.Enable()
	if ib.Disabled() {
		t.Fatal("expected badge to be re-enabled")
	}
}

func TestShadcnButton_VariantsAndActions(t *testing.T) {
	clicked := false
	variants := []ButtonVariant{
		ButtonDefault,
		ButtonSecondary,
		ButtonDestructive,
		ButtonOutline,
		ButtonGhost,
	}

	iconRes := ResourceFromSVG("test_btn.svg", SVGShieldCheck)

	for _, v := range variants {
		btn := NewShadcnButton("Action", v, ButtonSizeDefault, iconRes, func() {
			clicked = true
		})
		if btn == nil {
			t.Fatalf("expected non-nil button for variant %d", v)
		}
		btn.Tapped(nil)
		if !clicked {
			t.Fatal("expected button click handler to fire")
		}
		clicked = false
	}

	// Icon-only button
	iconOnlyBtn := NewShadcnButton("", ButtonDefault, ButtonSizeIcon, iconRes, nil)
	if iconOnlyBtn == nil {
		t.Fatal("expected non-nil icon-only button")
	}

	// Text-only button
	textOnlyBtn := NewShadcnButton("Simple", ButtonSecondary, ButtonSizeDefault, nil, nil)
	if textOnlyBtn == nil {
		t.Fatal("expected non-nil text-only button")
	}
}

func TestShadcnCard_CompoundStructure(t *testing.T) {
	badge := NewBadge("ACTIVE", BadgeSuccess, BadgeShapePill)
	content := widget.NewLabel("Body Content")
	footer := widget.NewLabel("Footer Metadata")

	card := NewShadcnCard(CardParts{
		Title:       "Card Title",
		Description: "Card Description Helper",
		Badge:       badge,
		Content:     content,
		Footer:      footer,
	})

	if card == nil {
		t.Fatal("expected non-nil compound card")
	}

	// Card without header or footer
	simpleCard := NewShadcnCard(CardParts{
		Content: content,
	})
	if simpleCard == nil {
		t.Fatal("expected non-nil simple card")
	}
}

func TestShadcnInput_AndFormField(t *testing.T) {
	// Standard input
	inp := NewShadcnInput("Placeholder text", false)
	if inp == nil || inp.PlaceHolder != "Placeholder text" {
		t.Fatal("expected valid standard input entry")
	}

	// Password input
	pwdInp := NewShadcnInput("Password", true)
	if pwdInp == nil || !pwdInp.Password {
		t.Fatal("expected valid password entry")
	}

	// FormField compound
	fieldObj, entry := NewFormField(FormField{
		Label:       "Full Name",
		Placeholder: "e.g. Yogesh Mallik",
		HelperText:  "As printed on academic marksheet",
		IsPassword:  false,
	})

	if fieldObj == nil || entry == nil {
		t.Fatal("expected non-nil form field and entry")
	}

	// FormField with pre-existing entry
	customField, _ := NewFormField(FormField{
		Label: "Existing",
		Entry: entry,
	})
	if customField == nil {
		t.Fatal("expected non-nil form field with custom entry")
	}

	// FormField disabled state
	disabledField, disEntry := NewFormField(FormField{
		Label:    "Roll Number",
		Disabled: true,
	})
	if disabledField == nil || disEntry == nil {
		t.Fatal("expected non-nil disabled form field")
	}
	if !disEntry.Disabled() {
		t.Fatal("expected entry to be disabled when FormField.Disabled is true")
	}

	// FormField with error state
	errorField, _ := NewFormField(FormField{
		Label:     "Token",
		ErrorText: "Token has expired or is invalid",
	})
	if errorField == nil {
		t.Fatal("expected non-nil error form field")
	}
}

func TestShadcnAlert_Variants(t *testing.T) {
	variants := []AlertVariant{
		AlertDefault,
		AlertDestructive,
		AlertWarning,
		AlertSuccess,
	}

	iconRes := ResourceFromSVG("test_alert.svg", SVGShieldCheck)

	for _, v := range variants {
		alertWithIcon := NewShadcnAlert("Alert Title", "Detailed description of system event", v, iconRes)
		if alertWithIcon == nil {
			t.Fatalf("expected non-nil alert with icon for variant %d", v)
		}

		alertNoIcon := NewShadcnAlert("Title Only", "Description without icon", v, nil)
		if alertNoIcon == nil {
			t.Fatalf("expected non-nil alert without icon for variant %d", v)
		}
	}
}
