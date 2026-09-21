package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// EmptyStateParams defines parts for rendering an empty state view.
type EmptyStateParams struct {
	Media       fyne.CanvasObject
	Title       string
	Description string
	Action      fyne.CanvasObject
}

// NewEmptyState constructs a centered empty state component matching shadcn empty.svelte.
func NewEmptyState(params EmptyStateParams) fyne.CanvasObject {
	items := make([]fyne.CanvasObject, 0, 4)

	if params.Media != nil {
		items = append(items, container.NewCenter(params.Media))
	}

	if params.Title != "" {
		titleLbl := widget.NewLabelWithStyle(params.Title, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
		titleLbl.Wrapping = fyne.TextWrapWord
		items = append(items, titleLbl)
	}

	if params.Description != "" {
		descLbl := widget.NewLabelWithStyle(params.Description, fyne.TextAlignCenter, fyne.TextStyle{})
		descLbl.Wrapping = fyne.TextWrapWord
		items = append(items, descLbl)
	}

	if params.Action != nil {
		items = append(items, container.NewCenter(params.Action))
	}

	box := container.NewVBox(items...)
	card := container.NewStack(
		NewShadcnCard(CardParts{Content: container.NewPadded(box)}),
	)

	_ = theme.ColorNameForeground // reference check
	return container.NewCenter(card)
}
