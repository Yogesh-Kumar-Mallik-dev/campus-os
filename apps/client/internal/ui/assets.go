package ui

import (
	_ "embed"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

//go:embed assets/bbdit-logo-transparent.png
var bbditLogoBytes []byte

// BBDITLogoResource encapsulates the official institutional logo image.
var BBDITLogoResource = fyne.NewStaticResource("bbdit-logo-transparent.png", bbditLogoBytes)

// RenderLogoImage renders the raw official BBDIT institutional logo with aspect ratio preserved.
func RenderLogoImage(width, height float32) *canvas.Image {
	img := canvas.NewImageFromResource(BBDITLogoResource)
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(width, height))
	return img
}

// RenderBBDITHeaderLogo renders the official BBDIT institutional logo within a crisp high-contrast pill
// with subtle border and tight padding, ensuring seamless blending across light and dark themes.
func RenderBBDITHeaderLogo(width, height float32) fyne.CanvasObject {
	if width <= 0 {
		width = 120
	}
	if height <= 0 {
		height = 24
	}

	logoImg := RenderLogoImage(width, height)
	pillBg := canvas.NewRectangle(color.NRGBA{R: 255, G: 255, B: 255, A: 235})
	pillBg.CornerRadius = 4
	pillBg.StrokeColor = color.NRGBA{R: 203, G: 213, B: 225, A: 160}
	pillBg.StrokeWidth = 1

	return container.NewStack(
		pillBg,
		container.NewCenter(logoImg),
	)
}
