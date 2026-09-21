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
// matching the exact presentation used in BBDIT's web/public layout (bg-white/95 with rounded corners).
func RenderBBDITHeaderLogo(width, height float32) fyne.CanvasObject {
	if width <= 0 {
		width = 160
	}
	if height <= 0 {
		height = 28
	}

	logoImg := RenderLogoImage(width, height)
	pillBg := canvas.NewRectangle(color.NRGBA{R: 255, G: 255, B: 255, A: 245})
	pillBg.CornerRadius = 6

	return container.NewStack(
		pillBg,
		container.NewPadded(logoImg),
	)
}
