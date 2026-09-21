package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// CardParts defines structural sections of a compound shadcn Card.
type CardParts struct {
	Title       string
	Description string
	Badge       fyne.CanvasObject
	Content     fyne.CanvasObject
	Footer      fyne.CanvasObject
}

// NewShadcnCard constructs a compound card widget mirroring shadcn-svelte's Card compound architecture.
func NewShadcnCard(parts CardParts) fyne.CanvasObject {
	cardBg := canvas.NewRectangle(theme.Color(theme.ColorNameMenuBackground))
	cardBg.CornerRadius = 10
	cardBg.StrokeColor = theme.Color(theme.ColorNameInputBorder)
	cardBg.StrokeWidth = 1

	elements := make([]fyne.CanvasObject, 0, 4)

	// Header section
	if parts.Title != "" || parts.Description != "" || parts.Badge != nil {
		headerItems := make([]fyne.CanvasObject, 0, 3)

		if parts.Badge != nil {
			headerItems = append(headerItems, container.NewHBox(parts.Badge))
		}

		if parts.Title != "" {
			titleLbl := widget.NewLabelWithStyle(parts.Title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
			titleLbl.Wrapping = fyne.TextWrapWord
			headerItems = append(headerItems, titleLbl)
		}

		if parts.Description != "" {
			descLbl := widget.NewLabel(parts.Description)
			descLbl.Wrapping = fyne.TextWrapWord
			headerItems = append(headerItems, descLbl)
		}

		headerItems = append(headerItems, widget.NewSeparator())
		elements = append(elements, container.NewVBox(headerItems...))
	}

	// Content section
	if parts.Content != nil {
		elements = append(elements, parts.Content)
	}

	// Footer section
	if parts.Footer != nil {
		elements = append(elements, widget.NewSeparator(), parts.Footer)
	}

	return container.NewStack(
		cardBg,
		container.NewPadded(container.NewVBox(elements...)),
	)
}
