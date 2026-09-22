package layout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"image/color"
)

// Spacer returns a flexible spacer item with weight 1.0 that expands to fill available space in Rows or Columns.
func Spacer() fyne.CanvasObject {
	rect := canvas.NewRectangle(color.Transparent)
	rect.SetMinSize(fyne.NewSize(0, 0))
	return FlexItem(rect, 1.0)
}

// FixedSpacer returns a transparent spacer with fixed minimum dimensions.
func FixedSpacer(width, height float32) fyne.CanvasObject {
	rect := canvas.NewRectangle(color.Transparent)
	rect.SetMinSize(fyne.NewSize(width, height))
	return rect
}
