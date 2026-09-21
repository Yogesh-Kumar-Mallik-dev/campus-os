package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

// BadgeVariant enumerates the shadcn badge style variants.
type BadgeVariant int

const (
	BadgeDefault BadgeVariant = iota
	BadgeSecondary
	BadgeDestructive
	BadgeOutline
	BadgeSuccess
	BadgeWarning
)

// BadgeShape controls the corner geometry of the badge.
type BadgeShape int

const (
	BadgeShapePill BadgeShape = iota
	BadgeShapeRounded
)

// NewBadge constructs a shadcn-styled badge component matching institutional design tokens.
func NewBadge(text string, variant BadgeVariant, shape BadgeShape) fyne.CanvasObject {
	var bg color.Color
	var fg color.Color
	var border color.Color

	switch variant {
	case BadgeDefault:
		// Primary institutional terracotta: #F45A51
		bg = color.NRGBA{R: 244, G: 90, B: 81, A: 255}
		fg = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		border = color.NRGBA{R: 244, G: 90, B: 81, A: 255}
	case BadgeSecondary:
		// Secondary surface: #282E3E
		bg = color.NRGBA{R: 40, G: 46, B: 62, A: 255}
		fg = color.NRGBA{R: 242, G: 246, B: 248, A: 255}
		border = color.NRGBA{R: 46, G: 56, B: 68, A: 255}
	case BadgeDestructive:
		// Destructive red tint: #451818 bg, #F87171 fg
		bg = color.NRGBA{R: 69, G: 24, B: 24, A: 230}
		fg = color.NRGBA{R: 248, G: 113, B: 113, A: 255}
		border = color.NRGBA{R: 120, G: 35, B: 35, A: 255}
	case BadgeOutline:
		// Transparent with subtle border
		bg = color.NRGBA{R: 16, G: 22, B: 28, A: 0}
		fg = color.NRGBA{R: 242, G: 246, B: 248, A: 255}
		border = color.NRGBA{R: 60, G: 72, B: 92, A: 255}
	case BadgeSuccess:
		// Emerald success: #10402E bg, #34D399 fg
		bg = color.NRGBA{R: 16, G: 64, B: 46, A: 230}
		fg = color.NRGBA{R: 52, G: 211, B: 153, A: 255}
		border = color.NRGBA{R: 20, G: 100, B: 70, A: 255}
	case BadgeWarning:
		// Amber warning: #3D2A0E bg, #FBBF24 fg
		bg = color.NRGBA{R: 61, G: 42, B: 14, A: 230}
		fg = color.NRGBA{R: 251, G: 191, B: 36, A: 255}
		border = color.NRGBA{R: 120, G: 85, B: 20, A: 255}
	}

	bgRect := canvas.NewRectangle(bg)
	bgRect.StrokeColor = border
	bgRect.StrokeWidth = 1

	if shape == BadgeShapePill {
		bgRect.CornerRadius = 12
	} else {
		bgRect.CornerRadius = 5
	}

	badgeText := canvas.NewText("  "+text+"  ", fg)
	badgeText.TextSize = 11
	badgeText.TextStyle = fyne.TextStyle{Bold: true}
	badgeText.Alignment = fyne.TextAlignCenter

	return container.NewStack(
		bgRect,
		container.NewPadded(badgeText),
	)
}
