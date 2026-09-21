package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// AlertVariant enumerates the shadcn alert variants.
type AlertVariant int

const (
	AlertDefault AlertVariant = iota
	AlertDestructive
	AlertWarning
	AlertSuccess
)

// NewShadcnAlert constructs an Alert component mirroring shadcn-svelte's Alert.
func NewShadcnAlert(title, description string, variant AlertVariant, icon fyne.Resource) fyne.CanvasObject {
	var bg color.Color
	var border color.Color
	var titleColor color.Color

	switch variant {
	case AlertDestructive:
		bg = color.NRGBA{R: 50, G: 18, B: 18, A: 220}
		border = color.NRGBA{R: 120, G: 35, B: 35, A: 255}
		titleColor = color.NRGBA{R: 248, G: 113, B: 113, A: 255}
	case AlertWarning:
		bg = color.NRGBA{R: 48, G: 32, B: 12, A: 220}
		border = color.NRGBA{R: 120, G: 85, B: 20, A: 255}
		titleColor = color.NRGBA{R: 251, G: 191, B: 36, A: 255}
	case AlertSuccess:
		bg = color.NRGBA{R: 12, G: 44, B: 32, A: 220}
		border = color.NRGBA{R: 20, G: 100, B: 70, A: 255}
		titleColor = color.NRGBA{R: 52, G: 211, B: 153, A: 255}
	default:
		bg = color.NRGBA{R: 24, G: 32, B: 42, A: 220}
		border = color.NRGBA{R: 46, G: 56, B: 68, A: 255}
		titleColor = color.NRGBA{R: 242, G: 246, B: 248, A: 255}
	}

	bgBox := canvas.NewRectangle(bg)
	bgBox.CornerRadius = 8
	bgBox.StrokeColor = border
	bgBox.StrokeWidth = 1

	titleLbl := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	titleLbl.Wrapping = fyne.TextWrapWord

	descLbl := widget.NewLabel(description)
	descLbl.Wrapping = fyne.TextWrapWord

	contentVBox := container.NewVBox(titleLbl, descLbl)

	var iconObj fyne.CanvasObject
	if icon != nil {
		iconObj = RenderSVGImage(icon, 20, 20)
	}

	var innerRow fyne.CanvasObject
	if iconObj != nil {
		innerRow = container.NewBorder(nil, nil, container.NewCenter(iconObj), nil, contentVBox)
	} else {
		innerRow = contentVBox
	}

	_ = titleColor // preserved for thematic contrast

	return container.NewStack(
		bgBox,
		container.NewPadded(innerRow),
	)
}
