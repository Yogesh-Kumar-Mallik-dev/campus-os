package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// CampusTheme implements fyne.Theme strictly mapped to the institutional design system
// defined in frontend/web/src/routes/+layout.css:
// Light mode:
//   --primary: oklch(0.42 0.16 27)    -> rgb(144, 20, 22)
//   --background: oklch(0.978 0.005 240) -> rgb(245, 248, 251)
//   --foreground: oklch(0.18 0.015 240)  -> rgb(16, 22, 28)
//   --card: oklch(1 0 0)                  -> rgb(255, 255, 255)
//   --border: oklch(0.895 0.008 240)      -> rgb(216, 221, 225)
// Dark mode:
//   --primary: oklch(0.67 0.19 27)    -> rgb(244, 90, 81)
//   --background: oklch(0.16 0.018 242) -> rgb(16, 22, 28)
//   --foreground: oklch(0.97 0.004 225)  -> rgb(242, 246, 248)
//   --card: oklch(0.205 0.02 240)        -> rgb(24, 32, 42)
//   --muted-foreground: oklch(0.72 0.018 225) -> rgb(153, 167, 173)
//   --border: oklch(0.72 0.018 225 / 16%)     -> rgb(38, 47, 58)
type CampusTheme struct {
	fyne.Theme
	forcedVariant *fyne.ThemeVariant
}

var (
	currentThemeVariant fyne.ThemeVariant = theme.VariantDark
)

// IsDarkTheme returns whether dark mode is currently active.
func IsDarkTheme() bool {
	return currentThemeVariant == theme.VariantDark
}

// ToggleTheme switches between light and dark institutional themes across the application.
func ToggleTheme() {
	if currentThemeVariant == theme.VariantDark {
		currentThemeVariant = theme.VariantLight
	} else {
		currentThemeVariant = theme.VariantDark
	}
	if a := fyne.CurrentApp(); a != nil && a.Settings() != nil {
		a.Settings().SetTheme(NewCampusThemeWithVariant(currentThemeVariant))
	}
}

// NewCampusTheme creates an institutional theme with colors from +layout.css.
func NewCampusTheme() fyne.Theme {
	return &CampusTheme{
		Theme: theme.DefaultTheme(),
	}
}

// NewCampusThemeWithVariant creates an institutional theme forced to a specific variant.
func NewCampusThemeWithVariant(variant fyne.ThemeVariant) fyne.Theme {
	return &CampusTheme{
		Theme:         theme.DefaultTheme(),
		forcedVariant: &variant,
	}
}

// Color returns tailored colors mapped from the institutional +layout.css.
func (t *CampusTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	activeVariant := variant
	if t.forcedVariant != nil {
		activeVariant = *t.forcedVariant
	}
	if activeVariant == theme.VariantLight {
		switch name {
		case theme.ColorNameBackground:
			// --background: oklch(0.978 0.005 240)
			return color.NRGBA{R: 245, G: 248, B: 251, A: 255}
		case theme.ColorNameForeground:
			// --foreground: oklch(0.18 0.015 240)
			return color.NRGBA{R: 16, G: 22, B: 28, A: 255}
		case theme.ColorNameMenuBackground, theme.ColorNameOverlayBackground, theme.ColorNameHeaderBackground:
			// --card: oklch(1 0 0)
			return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		case theme.ColorNameInputBackground:
			return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		case theme.ColorNameInputBorder:
			// --border: oklch(0.895 0.008 240)
			return color.NRGBA{R: 216, G: 221, B: 225, A: 255}
		case theme.ColorNamePrimary:
			// --primary: oklch(0.42 0.16 27)
			return color.NRGBA{R: 144, G: 20, B: 22, A: 255}
		case theme.ColorNameForegroundOnPrimary:
			return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		case theme.ColorNameButton:
			return color.NRGBA{R: 237, G: 240, B: 244, A: 255}
		case theme.ColorNameDisabledButton:
			return color.NRGBA{R: 230, G: 234, B: 240, A: 255}
		case theme.ColorNameHover:
			// Subtle darkening alpha blend overlay for light backgrounds
			return color.NRGBA{R: 0, G: 0, B: 0, A: 15}
		case theme.ColorNamePressed:
			return color.NRGBA{R: 0, G: 0, B: 0, A: 32}
		case theme.ColorNameFocus:
			return color.NRGBA{R: 144, G: 20, B: 22, A: 220}
		case theme.ColorNameSelection:
			return color.NRGBA{R: 144, G: 20, B: 22, A: 60}
		case theme.ColorNameHyperlink:
			return color.NRGBA{R: 144, G: 20, B: 22, A: 255}
		case theme.ColorNameSeparator:
			return color.NRGBA{R: 216, G: 221, B: 225, A: 255}
		case theme.ColorNamePlaceHolder:
			return color.NRGBA{R: 101, G: 108, B: 125, A: 255}
		case theme.ColorNameDisabled:
			return color.NRGBA{R: 130, G: 140, B: 155, A: 255}
		case theme.ColorNameSuccess:
			return color.NRGBA{R: 16, G: 185, B: 129, A: 255}
		case theme.ColorNameWarning:
			return color.NRGBA{R: 245, G: 158, B: 11, A: 255}
		case theme.ColorNameError:
			// --destructive: oklch(0.577 0.245 27.325)
			return color.NRGBA{R: 220, G: 38, B: 38, A: 255}
		case theme.ColorNameForegroundOnError:
			return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		case theme.ColorNameForegroundOnSuccess:
			return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		case theme.ColorNameForegroundOnWarning:
			return color.NRGBA{R: 16, G: 22, B: 28, A: 255}
		case theme.ColorNameScrollBar:
			return color.NRGBA{R: 180, G: 186, B: 196, A: 180}
		case theme.ColorNameScrollBarBackground:
			return color.NRGBA{R: 245, G: 248, B: 251, A: 0}
		case theme.ColorNameShadow:
			return color.NRGBA{R: 0, G: 0, B: 0, A: 30}
		case theme.ColorNameInnerWindowBorder:
			return color.NRGBA{R: 216, G: 221, B: 225, A: 255}
		case theme.ColorNameInnerWindowBorderInactive:
			return color.NRGBA{R: 235, G: 238, B: 242, A: 255}
		}
	} else {
		// Dark theme (institution's primary dark palette)
		switch name {
		case theme.ColorNameBackground:
			// --background: oklch(0.16 0.018 242)
			return color.NRGBA{R: 16, G: 22, B: 28, A: 255}
		case theme.ColorNameForeground:
			// --foreground: oklch(0.97 0.004 225)
			return color.NRGBA{R: 242, G: 246, B: 248, A: 255}
		case theme.ColorNameMenuBackground, theme.ColorNameOverlayBackground, theme.ColorNameHeaderBackground:
			// --card: oklch(0.205 0.02 240)
			return color.NRGBA{R: 24, G: 32, B: 42, A: 255}
		case theme.ColorNameInputBackground:
			return color.NRGBA{R: 19, G: 26, B: 35, A: 255}
		case theme.ColorNameInputBorder:
			// --border: oklch(0.72 0.018 225 / 16%)
			return color.NRGBA{R: 46, G: 56, B: 68, A: 255}
		case theme.ColorNamePrimary:
			// --primary: oklch(0.67 0.19 27) -> #f45a51
			return color.NRGBA{R: 244, G: 90, B: 81, A: 255}
		case theme.ColorNameForegroundOnPrimary:
			// Crisp pure white text and icons so buttons shine with high contrast against terracotta
			return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		case theme.ColorNameButton:
			// --secondary: oklch(0.27 0.018 239)
			return color.NRGBA{R: 40, G: 46, B: 62, A: 255}
		case theme.ColorNameDisabledButton:
			return color.NRGBA{R: 24, G: 30, B: 40, A: 255}
		case theme.ColorNameHover:
			// Subtle luminous highlight alpha blend overlay so button colors brighten rather than greying out
			return color.NRGBA{R: 255, G: 255, B: 255, A: 26}
		case theme.ColorNamePressed:
			return color.NRGBA{R: 255, G: 255, B: 255, A: 50}
		case theme.ColorNameFocus:
			return color.NRGBA{R: 244, G: 90, B: 81, A: 220}
		case theme.ColorNameSelection:
			return color.NRGBA{R: 244, G: 90, B: 81, A: 90}
		case theme.ColorNameHyperlink:
			return color.NRGBA{R: 244, G: 90, B: 81, A: 255}
		case theme.ColorNameSeparator:
			return color.NRGBA{R: 46, G: 56, B: 68, A: 255}
		case theme.ColorNamePlaceHolder:
			// --muted-foreground: oklch(0.72 0.018 225)
			return color.NRGBA{R: 153, G: 167, B: 173, A: 220}
		case theme.ColorNameDisabled:
			return color.NRGBA{R: 110, G: 122, B: 138, A: 255}
		case theme.ColorNameSuccess:
			return color.NRGBA{R: 52, G: 211, B: 153, A: 255}
		case theme.ColorNameWarning:
			return color.NRGBA{R: 251, G: 191, B: 36, A: 255}
		case theme.ColorNameError:
			// --destructive: oklch(0.704 0.191 22.216)
			return color.NRGBA{R: 248, G: 113, B: 113, A: 255}
		case theme.ColorNameForegroundOnError:
			return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		case theme.ColorNameForegroundOnSuccess:
			return color.NRGBA{R: 16, G: 22, B: 28, A: 255}
		case theme.ColorNameForegroundOnWarning:
			return color.NRGBA{R: 16, G: 22, B: 28, A: 255}
		case theme.ColorNameScrollBar:
			return color.NRGBA{R: 52, G: 62, B: 80, A: 180}
		case theme.ColorNameScrollBarBackground:
			return color.NRGBA{R: 16, G: 22, B: 28, A: 0}
		case theme.ColorNameShadow:
			return color.NRGBA{R: 0, G: 0, B: 0, A: 90}
		case theme.ColorNameInnerWindowBorder:
			return color.NRGBA{R: 46, G: 56, B: 68, A: 255}
		case theme.ColorNameInnerWindowBorderInactive:
			return color.NRGBA{R: 30, G: 38, B: 48, A: 255}
		}
	}

	return t.Theme.Color(name, variant)
}

// Size returns refined dimensions matching CSS tokens: --radius: 0.45rem.
func (t *CampusTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameCardRadius:
		return 10.0 // matching --radius
	case theme.SizeNameButtonRadius:
		return 8.0 // --radius-md
	case theme.SizeNameInputRadius:
		return 8.0
	case theme.SizeNameSelectionRadius:
		return 8.0
	case theme.SizeNamePadding:
		return 10.0
	case theme.SizeNameInnerPadding:
		return 8.0
	case theme.SizeNameText:
		return 13.5
	case theme.SizeNameHeadingText:
		return 19.0
	case theme.SizeNameSubHeadingText:
		return 15.0
	case theme.SizeNameCaptionText:
		return 11.0
	default:
		return t.Theme.Size(name)
	}
}
