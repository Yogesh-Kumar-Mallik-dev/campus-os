package ui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
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
	var bg, border, fg color.Color

	switch variant {
	case PillSuccess:
		bg = color.NRGBA{R: 16, G: 185, B: 129, A: 35}
		border = color.NRGBA{R: 16, G: 185, B: 129, A: 85}
		fg = color.NRGBA{R: 52, G: 211, B: 153, A: 255}
	case PillWarning:
		bg = color.NRGBA{R: 245, G: 158, B: 11, A: 35}
		border = color.NRGBA{R: 245, G: 158, B: 11, A: 85}
		fg = color.NRGBA{R: 251, G: 191, B: 36, A: 255}
	case PillError:
		bg = color.NRGBA{R: 239, G: 68, B: 68, A: 35}
		border = color.NRGBA{R: 239, G: 68, B: 68, A: 85}
		fg = color.NRGBA{R: 248, G: 113, B: 113, A: 255}
	case PillInfo:
		bg = color.NRGBA{R: 59, G: 130, B: 246, A: 35}
		border = color.NRGBA{R: 59, G: 130, B: 246, A: 85}
		fg = color.NRGBA{R: 96, G: 165, B: 250, A: 255}
	default:
		bg = color.NRGBA{R: 255, G: 255, B: 255, A: 14}
		border = color.NRGBA{R: 255, G: 255, B: 255, A: 28}
		fg = color.NRGBA{R: 226, G: 232, B: 240, A: 240}
	}

	pillBg := canvas.NewRectangle(bg)
	pillBg.CornerRadius = 12
	pillBg.StrokeColor = border
	pillBg.StrokeWidth = 1

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

// NewKeyValueRow renders a modern, scannable metadata row matching Linear/Raycast design standards.
// It pairs an uppercase muted micro-label with high-contrast text and an optional status badge.
func NewKeyValueRow(label, value string, badge fyne.CanvasObject) fyne.CanvasObject {
	lblText := canvas.NewText(label, color.NRGBA{R: 148, G: 163, B: 184, A: 220})
	lblText.TextSize = 11
	lblText.TextStyle = fyne.TextStyle{Bold: true}

	valLbl := widget.NewLabel(value)
	valLbl.Wrapping = fyne.TextWrapWord

	leftBox := container.NewVBox(lblText, valLbl)

	if badge != nil {
		return container.NewBorder(nil, nil, leftBox, container.NewCenter(badge))
	}
	return leftBox
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
	var iconBg color.Color
	var titleColor color.Color

	switch variant {
	case PillWarning:
		bg = color.NRGBA{R: 245, G: 158, B: 11, A: 22}
		border = color.NRGBA{R: 245, G: 158, B: 11, A: 65}
		iconBg = color.NRGBA{R: 245, G: 158, B: 11, A: 40}
		titleColor = color.NRGBA{R: 251, G: 191, B: 36, A: 255}
	case PillSuccess:
		bg = color.NRGBA{R: 16, G: 185, B: 129, A: 22}
		border = color.NRGBA{R: 16, G: 185, B: 129, A: 65}
		iconBg = color.NRGBA{R: 16, G: 185, B: 129, A: 40}
		titleColor = color.NRGBA{R: 52, G: 211, B: 153, A: 255}
	case PillError:
		bg = color.NRGBA{R: 239, G: 68, B: 68, A: 22}
		border = color.NRGBA{R: 239, G: 68, B: 68, A: 65}
		iconBg = color.NRGBA{R: 239, G: 68, B: 68, A: 40}
		titleColor = color.NRGBA{R: 248, G: 113, B: 113, A: 255}
	default:
		bg = color.NRGBA{R: 59, G: 130, B: 246, A: 20}
		border = color.NRGBA{R: 59, G: 130, B: 246, A: 60}
		iconBg = color.NRGBA{R: 59, G: 130, B: 246, A: 35}
		titleColor = color.NRGBA{R: 96, G: 165, B: 250, A: 255}
	}

	bgRect := canvas.NewRectangle(bg)
	bgRect.CornerRadius = 8
	bgRect.StrokeColor = border
	bgRect.StrokeWidth = 1

	titleText := canvas.NewText(title, titleColor)
	titleText.TextSize = 12
	titleText.TextStyle = fyne.TextStyle{Bold: true}

	msgText := widget.NewLabel(message)
	msgText.Wrapping = fyne.TextWrapWord

	textVBox := container.NewVBox(titleText, msgText)

	var iconBadge fyne.CanvasObject
	if svgIcon != nil {
		iconImg := RenderSVGImage(svgIcon, 16, 16)
		tileBg := canvas.NewRectangle(iconBg)
		tileBg.CornerRadius = 6
		tileBg.StrokeColor = border
		tileBg.StrokeWidth = 1

		badge := container.NewStack(tileBg, container.NewCenter(iconImg))
		iconBadge = container.NewVBox(container.NewGridWrap(fyne.NewSize(28, 28), badge))
	} else {
		iconImg := widget.NewIcon(theme.InfoIcon())
		iconBadge = container.NewVBox(container.NewGridWrap(fyne.NewSize(28, 28), container.NewCenter(iconImg)))
	}

	content := container.NewBorder(nil, nil, iconBadge, nil, textVBox)

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
		var stepBadge fyne.CanvasObject

		if stepNum < currentStep {
			// Completed step: secondary surface badge with checkmark
			stepBadge = NewBadge(fmt.Sprintf("%d ✓", stepNum), BadgeSecondary, BadgeShapePill)
		} else if stepNum == currentStep {
			// Active step: primary terracotta badge with step number & name
			stepBadge = NewBadge(fmt.Sprintf("%d: %s", stepNum, name), BadgeDefault, BadgeShapePill)
		} else {
			// Pending step: outline badge
			stepBadge = NewBadge(fmt.Sprintf("%d", stepNum), BadgeOutline, BadgeShapePill)
		}

		stepItems = append(stepItems, stepBadge)

		if i < len(stepNames)-1 {
			sep := canvas.NewLine(color.NRGBA{R: 60, G: 72, B: 92, A: 255})
			sep.StrokeWidth = 2
			stepItems = append(stepItems, container.NewCenter(sep))
		}
	}

	return container.NewHBox(stepItems...)
}

// NewTopBar creates the institutional top navigation branding bar with SVG logo and interactive theme toggle.
func NewTopBar(institutionName, userRole string) fyne.CanvasObject {
	bg := canvas.NewRectangle(theme.Color(theme.ColorNameMenuBackground))
	bg.StrokeColor = theme.Color(theme.ColorNameInputBorder)
	bg.StrokeWidth = 1

	logoPill := RenderBBDITHeaderLogo(140, 26)

	themeSwitch := NewSwitch(IsDarkTheme(), func(checked bool) {
		ToggleTheme()
	})
	themeRow := container.NewHBox(
		RenderSVGImage(ResourceFromSVG("sun.svg", LucideSun), 14, 14),
		themeSwitch,
		RenderSVGImage(ResourceFromSVG("moon.svg", LucideMoon), 14, 14),
	)

	var rightSide fyne.CanvasObject
	if userRole != "" {
		roleBadge := NewBadge(userRole, BadgeOutline, BadgeShapePill)
		rightSide = container.NewHBox(roleBadge, themeRow)
	} else {
		rightSide = themeRow
	}

	leftPad := canvas.NewRectangle(color.Transparent)
	leftPad.SetMinSize(fyne.NewSize(8, 1))

	barContent := container.NewBorder(
		nil, nil,
		container.NewHBox(leftPad, container.NewCenter(logoPill)),
		container.NewCenter(rightSide),
	)

	return container.NewStack(
		bg,
		container.NewPadded(barContent),
	)
}
