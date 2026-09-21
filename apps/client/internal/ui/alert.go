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

type alertLayout struct {
	padX float32
	padY float32
}

func (l *alertLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(objects) == 0 {
		return
	}
	// Object 0: background rectangle
	objects[0].Resize(size)
	objects[0].Move(fyne.NewPos(0, 0))

	// Object 1: content container
	if len(objects) > 1 {
		w := size.Width - (l.padX * 2)
		h := size.Height - (l.padY * 2)
		if w < 0 {
			w = 0
		}
		if h < 0 {
			h = 0
		}
		objects[1].Resize(fyne.NewSize(w, h))
		objects[1].Move(fyne.NewPos(l.padX, l.padY))
	}
}

func (l *alertLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if len(objects) < 2 {
		return fyne.NewSize(120, 36)
	}
	min := objects[1].MinSize()
	return fyne.NewSize(min.Width+(l.padX*2), min.Height+(l.padY*2))
}

// NewShadcnAlert constructs a modern translucent glass Alert component matching 2026 design standards.
func NewShadcnAlert(title, description string, variant AlertVariant, icon fyne.Resource) fyne.CanvasObject {
	var bg, border, iconBg, titleColor color.Color

	switch variant {
	case AlertDestructive:
		// Luminous Rose tint
		bg = color.NRGBA{R: 239, G: 68, B: 68, A: 22}
		border = color.NRGBA{R: 239, G: 68, B: 68, A: 65}
		iconBg = color.NRGBA{R: 239, G: 68, B: 68, A: 40}
		titleColor = color.NRGBA{R: 248, G: 113, B: 113, A: 255}
	case AlertWarning:
		// Warm Luminous Amber tint
		bg = color.NRGBA{R: 245, G: 158, B: 11, A: 22}
		border = color.NRGBA{R: 245, G: 158, B: 11, A: 65}
		iconBg = color.NRGBA{R: 245, G: 158, B: 11, A: 40}
		titleColor = color.NRGBA{R: 251, G: 191, B: 36, A: 255}
	case AlertSuccess:
		// Crisp Luminous Emerald tint
		bg = color.NRGBA{R: 16, G: 185, B: 129, A: 22}
		border = color.NRGBA{R: 16, G: 185, B: 129, A: 65}
		iconBg = color.NRGBA{R: 16, G: 185, B: 129, A: 40}
		titleColor = color.NRGBA{R: 52, G: 211, B: 153, A: 255}
	default:
		// Neutral Slate / Info tint
		bg = color.NRGBA{R: 59, G: 130, B: 246, A: 20}
		border = color.NRGBA{R: 59, G: 130, B: 246, A: 60}
		iconBg = color.NRGBA{R: 59, G: 130, B: 246, A: 35}
		titleColor = color.NRGBA{R: 96, G: 165, B: 250, A: 255}
	}

	bgBox := canvas.NewRectangle(bg)
	bgBox.CornerRadius = 8
	bgBox.StrokeColor = border
	bgBox.StrokeWidth = 1

	textItems := make([]fyne.CanvasObject, 0, 2)

	if title != "" {
		titleText := canvas.NewText(title, titleColor)
		titleText.TextSize = 12
		titleText.TextStyle = fyne.TextStyle{Bold: true}
		textItems = append(textItems, titleText)
	}

	if description != "" {
		descLbl := widget.NewLabel(description)
		descLbl.Wrapping = fyne.TextWrapWord
		textItems = append(textItems, descLbl)
	}

	textBox := container.NewVBox(textItems...)

	var content fyne.CanvasObject
	if icon != nil {
		iconImg := RenderSVGImage(icon, 16, 16)
		iconTileBg := canvas.NewRectangle(iconBg)
		iconTileBg.CornerRadius = 6
		iconTileBg.StrokeColor = border
		iconTileBg.StrokeWidth = 1

		iconBadge := container.NewStack(
			iconTileBg,
			container.NewCenter(iconImg),
		)
		iconBadgeWrapper := container.NewVBox(
			container.NewGridWrap(fyne.NewSize(28, 28), iconBadge),
		)

		content = container.NewBorder(nil, nil, iconBadgeWrapper, nil, textBox)
	} else {
		content = textBox
	}

	return container.New(&alertLayout{padX: 12, padY: 9}, bgBox, content)
}
