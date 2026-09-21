package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

// responsiveCardLayout dynamically centers content on wide desktop screens
// and fluidly expands to 100% width with touch-friendly margins on mobile viewports.
type responsiveCardLayout struct {
	maxWidth float32
	margin   float32
}

func NewResponsiveLayout(maxWidth, margin float32) fyne.Layout {
	if maxWidth <= 0 {
		maxWidth = 640
	}
	if margin <= 0 {
		margin = 12
	}
	return &responsiveCardLayout{
		maxWidth: maxWidth,
		margin:   margin,
	}
}

func (r *responsiveCardLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, child := range objects {
		targetW := size.Width - (r.margin * 2)
		if targetW > r.maxWidth {
			targetW = r.maxWidth
		}
		if targetW < 240 {
			targetW = 240
		}

		posX := (size.Width - targetW) / 2
		if posX < r.margin {
			posX = r.margin
		}

		minH := child.MinSize().Height
		targetH := size.Height
		if targetH < minH {
			targetH = minH
		}

		child.Move(fyne.NewPos(posX, 0))
		child.Resize(fyne.NewSize(targetW, targetH))
	}
}

func (r *responsiveCardLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	var maxMinW float32 = 260
	var maxMinH float32 = 0

	for _, child := range objects {
		ms := child.MinSize()
		if ms.Width > maxMinW {
			maxMinW = ms.Width
		}
		if ms.Height > maxMinH {
			maxMinH = ms.Height
		}
	}

	return fyne.NewSize(maxMinW, maxMinH)
}

// NewResponsiveCardContainer wraps a CanvasObject in a fluidly adapting container.
func NewResponsiveCardContainer(content fyne.CanvasObject, maxWidth float32) fyne.CanvasObject {
	inner := container.New(NewResponsiveLayout(maxWidth, 12), content)
	return container.NewVScroll(inner)
}
