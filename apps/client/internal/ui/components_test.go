package ui

import (
	"testing"

	"fyne.io/fyne/v2/theme"
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

	// Size checks
	cardRadius := th.Size(theme.SizeNameCardRadius)
	if cardRadius != 12.0 {
		t.Errorf("expected card radius 12.0, got %f", cardRadius)
	}

	btnRadius := th.Size(theme.SizeNameButtonRadius)
	if btnRadius != 8.0 {
		t.Errorf("expected button radius 8.0, got %f", btnRadius)
	}

	headingSize := th.Size(theme.SizeNameHeadingText)
	if headingSize != 20.0 {
		t.Errorf("expected heading size 20.0, got %f", headingSize)
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

	// Stat Card
	stat := NewStatCard("Metric", "100%", "Subtext", PillSuccess)
	if stat == nil {
		t.Fatal("expected non-nil stat card")
	}

	// Styled Card
	c1 := NewStyledCard("Card Title", stat)
	if c1 == nil {
		t.Fatal("expected non-nil styled card with title")
	}

	c2 := NewStyledCard("", stat)
	if c2 == nil {
		t.Fatal("expected non-nil styled card without title")
	}

	// Top Bar
	topBar := NewTopBar("CAMPUS OS", "Student Scholar")
	if topBar == nil {
		t.Fatal("expected non-nil top bar")
	}
}
