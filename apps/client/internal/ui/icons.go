package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// SVG vector definitions for crisp rendering at any DPI without emojis or PNGs.
const (
	SVGShieldCheck = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#60a5fa" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
		<path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
		<path d="m9 12 2 2 4-4"/>
	</svg>`

	SVGQRCodeFrame = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 48 48" fill="none" stroke="#77baff" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
		<path d="M6 16V8a2 2 0 0 1 2-2h8"/>
		<path d="M32 6h8a2 2 0 0 1 2 2v8"/>
		<path d="M42 32v8a2 2 0 0 1-2 2h-8"/>
		<path d="M16 42H8a2 2 0 0 1-2-2v-8"/>
		<rect x="12" y="12" width="8" height="8" rx="1.5" stroke="#3b82f6" fill="#1e3a5f"/>
		<rect x="28" y="12" width="8" height="8" rx="1.5" stroke="#3b82f6" fill="#1e3a5f"/>
		<rect x="12" y="28" width="8" height="8" rx="1.5" stroke="#3b82f6" fill="#1e3a5f"/>
		<circle cx="32" cy="32" r="2.5" fill="#60a5fa"/>
		<path d="M10 24h28" stroke="#38bdf8" stroke-width="1.5" stroke-dasharray="2 2"/>
	</svg>`

	SVGSIMCardActive = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#10b981" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
		<path d="M6 2h9l5 5v13a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2z"/>
		<path d="M9 11h6"/>
		<path d="M9 15h6"/>
		<path d="M9 7v1"/>
		<circle cx="15" cy="18" r="1" fill="#10b981"/>
	</svg>`

	SVGSIMCardSecondary = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#64748b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
		<path d="M6 2h9l5 5v13a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2z"/>
		<path d="M9 11h6"/>
		<path d="M9 15h6"/>
		<path d="M9 7v1"/>
	</svg>`

	SVGFingerprint = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#38bdf8" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
		<path d="M12 10a2 2 0 0 0-2 2c0 1.02-.1 2.51-.26 4"/>
		<path d="M14 13.12c0 2.38 0 6.38-1 8.88"/>
		<path d="M2 12h1"/>
		<path d="M21 12h1"/>
		<path d="M5 19.5C5.5 18 6 15 6 12a6 6 0 0 1 .34-2"/>
		<path d="M8.65 22c.21-.66.45-1.32.57-2"/>
		<path d="M9 6.8a6 6 0 0 1 9 5.2v2"/>
		<path d="M14 3.5c2.32.74 4 2.92 4 5.5v1"/>
	</svg>`

	SVGClockGrace = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#f59e0b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
		<circle cx="12" cy="12" r="10"/>
		<polyline points="12 6 12 12 16 14"/>
	</svg>`

	SVGAcademicCap = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#38bdf8" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
		<path d="M22 10v6M2 10l10-5 10 5-10 5z"/>
		<path d="M6 12v5c0 2 2 3 6 3s6-1 6-3v-5"/>
	</svg>`

	SVGCheckVerified = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#10b981" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
		<path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
		<polyline points="22 4 12 14.01 9 11.01"/>
	</svg>`

	SVGCampusLogo = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32" fill="none">
		<rect width="32" height="32" rx="8" fill="#1e293b"/>
		<path d="M16 6L6 11L16 16L26 11L16 6Z" fill="#38bdf8"/>
		<path d="M9 14.5V20.5C9 23.5 12 26 16 26C20 26 23 23.5 23 20.5V14.5L16 18L9 14.5Z" fill="#2563eb"/>
		<circle cx="16" cy="11" r="2" fill="#ffffff"/>
	</svg>`
)

// Resource helpers
func ResourceFromSVG(name, svgContent string) fyne.Resource {
	return fyne.NewStaticResource(name, []byte(svgContent))
}

// RenderSVGImage creates an SVG canvas image with explicit bounding dimensions.
func RenderSVGImage(res fyne.Resource, width, height float32) *canvas.Image {
	img := canvas.NewImageFromResource(res)
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(width, height))
	return img
}
