package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type PillVariant int

const (
	PillSuccess PillVariant = iota
	PillWarning
	PillError
	PillInfo
	PillNeutral
)

// NewStatusPill returns a custom styled badge capsule for status values.
func NewStatusPill(text string, variant PillVariant) fyne.CanvasObject {
	var bg color.Color
	var fg color.Color

	switch variant {
	case PillSuccess:
		bg = color.NRGBA{R: 16, G: 80, B: 55, A: 220}
		fg = color.NRGBA{R: 52, G: 211, B: 153, A: 255}
	case PillWarning:
		bg = color.NRGBA{R: 90, G: 65, B: 15, A: 220}
		fg = color.NRGBA{R: 251, G: 191, B: 36, A: 255}
	case PillError:
		bg = color.NRGBA{R: 90, G: 25, B: 25, A: 220}
		fg = color.NRGBA{R: 248, G: 113, B: 113, A: 255}
	case PillInfo:
		bg = color.NRGBA{R: 20, G: 55, B: 95, A: 220}
		fg = color.NRGBA{R: 119, G: 186, B: 255, A: 255}
	default:
		bg = color.NRGBA{R: 35, G: 45, B: 65, A: 220}
		fg = color.NRGBA{R: 200, G: 210, B: 230, A: 255}
	}

	pillBg := canvas.NewRectangle(bg)
	pillBg.CornerRadius = 12

	pillText := canvas.NewText("  "+text+"  ", fg)
	pillText.TextSize = 11
	pillText.TextStyle = fyne.TextStyle{Bold: true}
	pillText.Alignment = fyne.TextAlignCenter

	return container.NewStack(
		pillBg,
		container.NewPadded(pillText),
	)
}

// NewPageHeader constructs a styled domain screen banner.
func NewPageHeader(title, subtitle string, pill fyne.CanvasObject) fyne.CanvasObject {
	titleText := canvas.NewText(title, theme.Color(theme.ColorNameForeground))
	titleText.TextSize = 20
	titleText.TextStyle = fyne.TextStyle{Bold: true}

	subText := canvas.NewText(subtitle, theme.Color(theme.ColorNamePlaceHolder))
	subText.TextSize = 13

	leftBox := container.NewVBox(titleText, subText)

	if pill != nil {
		return container.NewBorder(nil, nil, leftBox, pill)
	}
	return leftBox
}

// NewStatCard builds a KPI / Metric widget box.
func NewStatCard(label, value string, sublabel string, pillVariant PillVariant) fyne.CanvasObject {
	lblText := canvas.NewText(label, theme.Color(theme.ColorNamePlaceHolder))
	lblText.TextSize = 12

	valText := canvas.NewText(value, theme.Color(theme.ColorNameForeground))
	valText.TextSize = 24
	valText.TextStyle = fyne.TextStyle{Bold: true}

	subText := canvas.NewText(sublabel, theme.Color(theme.ColorNamePlaceHolder))
	subText.TextSize = 11

	cardBg := canvas.NewRectangle(theme.Color(theme.ColorNameMenuBackground))
	cardBg.CornerRadius = 12
	cardBg.StrokeColor = theme.Color(theme.ColorNameInputBorder)
	cardBg.StrokeWidth = 1

	content := container.NewVBox(
		lblText,
		valText,
		subText,
	)

	return container.NewStack(
		cardBg,
		container.NewPadded(content),
	)
}

// NewStyledCard wraps any content into a card with uniform padding and rounded corners.
func NewStyledCard(title string, content fyne.CanvasObject) fyne.CanvasObject {
	cardBg := canvas.NewRectangle(theme.Color(theme.ColorNameMenuBackground))
	cardBg.CornerRadius = 12
	cardBg.StrokeColor = theme.Color(theme.ColorNameInputBorder)
	cardBg.StrokeWidth = 1

	var header fyne.CanvasObject
	if title != "" {
		headerText := canvas.NewText(title, theme.Color(theme.ColorNameForeground))
		headerText.TextSize = 15
		headerText.TextStyle = fyne.TextStyle{Bold: true}
		header = container.NewVBox(headerText, widget.NewSeparator())
	}

	body := content
	if header != nil {
		body = container.NewVBox(header, content)
	}

	return container.NewStack(
		cardBg,
		container.NewPadded(body),
	)
}

// NewTopBar creates the institutional top navigation branding bar.
func NewTopBar(institutionName, userRole string) fyne.CanvasObject {
	bg := canvas.NewRectangle(theme.Color(theme.ColorNameMenuBackground))
	bg.StrokeColor = theme.Color(theme.ColorNameInputBorder)
	bg.StrokeWidth = 1

	brandText := canvas.NewText(institutionName, theme.Color(theme.ColorNamePrimary))
	brandText.TextSize = 16
	brandText.TextStyle = fyne.TextStyle{Bold: true}

	tagline := canvas.NewText("Institutional Operating System", theme.Color(theme.ColorNamePlaceHolder))
	tagline.TextSize = 11

	brandBox := container.NewVBox(brandText, tagline)

	rolePill := NewStatusPill(userRole, PillInfo)
	statusDot := canvas.NewCircle(color.NRGBA{R: 52, G: 211, B: 153, A: 255})
	statusDot.Resize(fyne.NewSize(8, 8))
	liveLabel := canvas.NewText("Online", theme.Color(theme.ColorNamePlaceHolder))
	liveLabel.TextSize = 11

	rightBox := container.NewHBox(
		container.NewCenter(statusDot),
		liveLabel,
		layout.NewSpacer(),
		rolePill,
	)

	barContent := container.NewBorder(nil, nil, brandBox, rightBox)
	return container.NewStack(
		bg,
		container.NewPadded(barContent),
	)
}
