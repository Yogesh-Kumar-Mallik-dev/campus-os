package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

// NewShadcnSeparator creates a subtle divider line adhering to institutional border tokens.
func NewShadcnSeparator(horizontal bool) fyne.CanvasObject {
	rect := canvas.NewRectangle(theme.Color(theme.ColorNameSeparator))
	if horizontal {
		rect.SetMinSize(fyne.NewSize(1, 1))
	} else {
		rect.SetMinSize(fyne.NewSize(1, 1))
	}
	return rect
}
