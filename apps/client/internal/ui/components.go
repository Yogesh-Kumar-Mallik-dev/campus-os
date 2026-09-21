package ui

import (
	"fmt"
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

// NewPageHeader constructs a mobile-first responsive screen banner with word wrapping.
func NewPageHeader(title, subtitle string, pill fyne.CanvasObject) fyne.CanvasObject {
	titleLabel := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	titleLabel.Wrapping = fyne.TextWrapWord

	items := make([]fyne.CanvasObject, 0, 3)
	if pill != nil {
		items = append(items, container.NewHBox(pill))
	}
	items = append(items, titleLabel)

	if subtitle != "" {
		subLabel := widget.NewLabel(subtitle)
		subLabel.Wrapping = fyne.TextWrapWord
		items = append(items, subLabel)
	}

	return container.NewVBox(items...)
}

// NewStyledCard wraps any content into a card with uniform padding and rounded corners.
func NewStyledCard(title string, content fyne.CanvasObject) fyne.CanvasObject {
	cardBg := canvas.NewRectangle(theme.Color(theme.ColorNameMenuBackground))
	cardBg.CornerRadius = 10
	cardBg.StrokeColor = theme.Color(theme.ColorNameInputBorder)
	cardBg.StrokeWidth = 1

	var header fyne.CanvasObject
	if title != "" {
		headerText := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		headerText.Wrapping = fyne.TextWrapWord
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

// NewBannerNotice constructs a sleek notice banner with a vector SVG icon and subtle tint.
func NewBannerNotice(title, message string, variant PillVariant, svgIcon fyne.Resource) fyne.CanvasObject {
	var bg color.Color
	var border color.Color

	switch variant {
	case PillWarning:
		bg = color.NRGBA{R: 45, G: 32, B: 10, A: 200}
		border = color.NRGBA{R: 120, G: 85, B: 20, A: 255}
	case PillSuccess:
		bg = color.NRGBA{R: 10, G: 40, B: 30, A: 200}
		border = color.NRGBA{R: 20, G: 100, B: 70, A: 255}
	case PillError:
		bg = color.NRGBA{R: 45, G: 15, B: 15, A: 200}
		border = color.NRGBA{R: 120, G: 35, B: 35, A: 255}
	default:
		bg = color.NRGBA{R: 15, G: 30, B: 55, A: 200}
		border = color.NRGBA{R: 35, G: 75, B: 130, A: 255}
	}

	bgRect := canvas.NewRectangle(bg)
	bgRect.CornerRadius = 8
	bgRect.StrokeColor = border
	bgRect.StrokeWidth = 1

	titleText := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	titleText.Wrapping = fyne.TextWrapWord

	msgText := widget.NewLabel(message)
	msgText.Wrapping = fyne.TextWrapWord

	textVBox := container.NewVBox(titleText, msgText)

	var leftIcon fyne.CanvasObject
	if svgIcon != nil {
		leftIcon = RenderSVGImage(svgIcon, 24, 24)
	} else {
		leftIcon = widget.NewIcon(theme.InfoIcon())
	}

	content := container.NewBorder(nil, nil, container.NewCenter(leftIcon), nil, textVBox)

	return container.NewStack(
		bgRect,
		container.NewPadded(content),
	)
}

// NewStepIndicator renders a responsive multi-step progression bar that fits mobile & desktop.
func NewStepIndicator(currentStep int, stepNames []string) fyne.CanvasObject {
	stepItems := make([]fyne.CanvasObject, 0, len(stepNames)*2)

	for i, name := range stepNames {
		stepNum := i + 1
		var stepPill fyne.CanvasObject

		if stepNum < currentStep {
			// Completed step: compact checkmark number
			stepPill = NewStatusPill(fmt.Sprintf("%d", stepNum), PillSuccess)
		} else if stepNum == currentStep {
			// Active step: highlighted with full name
			stepPill = NewStatusPill(fmt.Sprintf("%d: %s", stepNum, name), PillInfo)
		} else {
			// Pending step: compact number
			stepPill = NewStatusPill(fmt.Sprintf("%d", stepNum), PillNeutral)
		}

		stepItems = append(stepItems, stepPill)

		if i < len(stepNames)-1 {
			sep := canvas.NewLine(color.NRGBA{R: 60, G: 72, B: 92, A: 255})
			sep.StrokeWidth = 2
			stepItems = append(stepItems, container.NewCenter(sep))
		}
	}

	return container.NewHBox(stepItems...)
}

// NewTopBar creates the institutional top navigation branding bar with SVG logo.
func NewTopBar(institutionName, userRole string) fyne.CanvasObject {
	bg := canvas.NewRectangle(theme.Color(theme.ColorNameMenuBackground))
	bg.StrokeColor = theme.Color(theme.ColorNameInputBorder)
	bg.StrokeWidth = 1

	logoPill := RenderBBDITHeaderLogo(120, 22)

	brandText := canvas.NewText(institutionName, theme.Color(theme.ColorNamePrimary))
	brandText.TextSize = 15
	brandText.TextStyle = fyne.TextStyle{Bold: true}

	tagline := canvas.NewText("Institutional Operating System", theme.Color(theme.ColorNamePlaceHolder))
	tagline.TextSize = 11

	brandBox := container.NewHBox(
		container.NewCenter(logoPill),
		container.NewVBox(brandText, tagline),
	)

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
