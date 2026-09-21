package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// CampusTheme implements fyne.Theme matching the design tokens from the abandoned project:
// --campus-color-surface, --campus-color-accent, --campus-color-text, --campus-color-muted
// and dark/light mode palette from tokens.css & +layout.css.
type CampusTheme struct {
	fyne.Theme
}

// NewCampusTheme creates an institutional theme with modern palette and geometry.
func NewCampusTheme() fyne.Theme {
	return &CampusTheme{
		Theme: theme.DefaultTheme(),
	}
}

// Color returns tailored colors mapped from design tokens.
func (t *CampusTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if variant == theme.VariantLight {
		switch name {
		case theme.ColorNameBackground:
			// --background: #f4f6f8 (tokens.css surface light)
			return color.NRGBA{R: 244, G: 246, B: 248, A: 255}
		case theme.ColorNameForeground:
			// --campus-color-text: #172033
			return color.NRGBA{R: 23, G: 32, B: 51, A: 255}
		case theme.ColorNameMenuBackground, theme.ColorNameOverlayBackground, theme.ColorNameHeaderBackground:
			return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		case theme.ColorNameInputBackground:
			return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		case theme.ColorNameInputBorder:
			return color.NRGBA{R: 226, G: 232, B: 240, A: 255}
		case theme.ColorNamePrimary:
			// --campus-color-accent (light): #075eb8
			return color.NRGBA{R: 7, G: 94, B: 184, A: 255}
		case theme.ColorNameForegroundOnPrimary:
			return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		case theme.ColorNameButton:
			return color.NRGBA{R: 235, G: 242, B: 250, A: 255}
		case theme.ColorNameHover:
			return color.NRGBA{R: 224, G: 234, B: 248, A: 255}
		case theme.ColorNamePressed:
			return color.NRGBA{R: 210, G: 225, B: 245, A: 255}
		case theme.ColorNameSeparator:
			return color.NRGBA{R: 226, G: 232, B: 240, A: 255}
		case theme.ColorNamePlaceHolder, theme.ColorNameDisabled:
			// --campus-color-muted: #56627a
			return color.NRGBA{R: 86, G: 98, B: 122, A: 255}
		case theme.ColorNameSuccess:
			return color.NRGBA{R: 16, G: 185, B: 129, A: 255}
		case theme.ColorNameWarning:
			return color.NRGBA{R: 245, G: 158, B: 11, A: 255}
		case theme.ColorNameError:
			return color.NRGBA{R: 239, G: 68, B: 68, A: 255}
		}
	} else {
		// Dark theme (default in modern workstations / Hyprland)
		switch name {
		case theme.ColorNameBackground:
			// --campus-color-surface (dark): rgba(16, 21, 34, 0.85) -> #101522
			return color.NRGBA{R: 16, G: 21, B: 34, A: 255}
		case theme.ColorNameForeground:
			// --campus-color-text: #f5f7fb
			return color.NRGBA{R: 245, G: 247, B: 251, A: 255}
		case theme.ColorNameMenuBackground, theme.ColorNameOverlayBackground, theme.ColorNameHeaderBackground:
			// --card: #1a2234
			return color.NRGBA{R: 26, G: 34, B: 52, A: 255}
		case theme.ColorNameInputBackground:
			// elevated subtle background
			return color.NRGBA{R: 21, G: 28, B: 44, A: 255}
		case theme.ColorNameInputBorder:
			// --border: rgba(180, 191, 211, 0.16) -> #2c3852
			return color.NRGBA{R: 44, G: 56, B: 82, A: 255}
		case theme.ColorNamePrimary:
			// --campus-color-accent (dark): #77baff
			return color.NRGBA{R: 119, G: 186, B: 255, A: 255}
		case theme.ColorNameForegroundOnPrimary:
			return color.NRGBA{R: 14, G: 20, B: 32, A: 255}
		case theme.ColorNameButton:
			return color.NRGBA{R: 30, G: 40, B: 62, A: 255}
		case theme.ColorNameHover:
			return color.NRGBA{R: 42, G: 56, B: 86, A: 255}
		case theme.ColorNamePressed:
			return color.NRGBA{R: 54, G: 72, B: 110, A: 255}
		case theme.ColorNameSeparator:
			return color.NRGBA{R: 38, G: 48, B: 70, A: 255}
		case theme.ColorNamePlaceHolder, theme.ColorNameDisabled:
			// --campus-color-muted: #b4bfd3
			return color.NRGBA{R: 180, G: 191, B: 211, A: 200}
		case theme.ColorNameSuccess:
			return color.NRGBA{R: 52, G: 211, B: 153, A: 255}
		case theme.ColorNameWarning:
			return color.NRGBA{R: 251, G: 191, B: 36, A: 255}
		case theme.ColorNameError:
			return color.NRGBA{R: 248, G: 113, B: 113, A: 255}
		}
	}

	return t.Theme.Color(name, variant)
}

// Size returns refined dimensions for modern desktop and mobile viewports.
func (t *CampusTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameCardRadius:
		return 12.0 // modern smooth card radius
	case theme.SizeNameButtonRadius:
		return 8.0 // sleek action button radius
	case theme.SizeNameInputRadius:
		return 8.0
	case theme.SizeNameSelectionRadius:
		return 8.0
	case theme.SizeNamePadding:
		return 10.0
	case theme.SizeNameInnerPadding:
		return 8.0
	case theme.SizeNameText:
		return 14.0
	case theme.SizeNameHeadingText:
		return 20.0
	case theme.SizeNameSubHeadingText:
		return 16.0
	case theme.SizeNameCaptionText:
		return 11.5
	default:
		return t.Theme.Size(name)
	}
}
