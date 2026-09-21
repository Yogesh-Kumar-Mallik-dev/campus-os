package ui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// SVG vector definitions mapped to the institutional color palette from +layout.css.
const (
	// SVGBBDITLogo represents the official BBDIT institutional crest with Vidya Sarvasya Bhushanam motto.
	SVGBBDITLogo = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 40 40" fill="none">
		<circle cx="20" cy="20" r="18" fill="#18202A" stroke="#f45a51" stroke-width="1.8"/>
		<circle cx="20" cy="20" r="14.5" stroke="#2E3844" stroke-width="1" stroke-dasharray="2 2"/>
		<!-- Outer Cogwheel teeth -->
		<path d="M20 7v2M20 31v2M7 20h2M31 20h2M11 11l1.5 1.5M27.5 27.5L29 29M11 29l1.5-1.5M27.5 12.5L29 11" stroke="#f45a51" stroke-width="1.5" stroke-linecap="round"/>
		<!-- Open Book of Learning -->
		<path d="M13 23c2.5-1 4.5-.5 7 1 2.5-1.5 4.5-2 7-1v-7c-2.5-1-4.5-.5-7 1-2.5-1.5-4.5-2-7-1v7z" fill="#222C38" stroke="#f2f6f8" stroke-width="1.2" stroke-linejoin="round"/>
		<!-- Science Atom Nucleus -->
		<ellipse cx="20" cy="18" rx="4.5" ry="1.8" stroke="#f45a51" stroke-width="0.9" transform="rotate(-25 20 18)"/>
		<ellipse cx="20" cy="18" rx="4.5" ry="1.8" stroke="#f45a51" stroke-width="0.9" transform="rotate(25 20 18)"/>
		<circle cx="20" cy="18" r="1" fill="#f45a51"/>
		<!-- Motto Banner Ribbon -->
		<path d="M10 28.5h20l-2 3H12l-2-3z" fill="#f45a51"/>
		<path d="M14 30h12" stroke="#ffffff" stroke-width="0.8" stroke-linecap="round"/>
	</svg>`

	// Lucide Icons (stroke: currentColor/brand, stroke-width: 2, linecap/linejoin: round)
	LucideScan = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#f45a51" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
		<path d="M3 7V5a2 2 0 0 1 2-2h2"/>
		<path d="M17 3h2a2 2 0 0 1 2 2v2"/>
		<path d="M21 17v2a2 2 0 0 1-2 2h-2"/>
		<path d="M7 21H5a2 2 0 0 1-2-2v-2"/>
	</svg>`

	LucideQrCode = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#f45a51" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
		<rect width="5" height="5" x="3" y="3" rx="1"/>
		<rect width="5" height="5" x="16" y="3" rx="1"/>
		<rect width="5" height="5" x="3" y="16" rx="1"/>
		<path d="M21 16h-3a2 2 0 0 0-2 2v3"/>
		<path d="M21 21v.01"/>
		<path d="M12 7v3a2 2 0 0 1-2 2H7"/>
		<path d="M3 12h.01"/>
		<path d="M12 3h.01"/>
		<path d="M12 16v.01"/>
		<path d="M16 12h1"/>
		<path d="M21 12v.01"/>
		<path d="M12 21v-1"/>
	</svg>`

	LucideSim = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#10b981" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
		<path d="M6 2h8l6 6v12a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2z"/>
		<path d="M10 10v4"/>
		<path d="M14 10v4"/>
		<path d="M10 14h4"/>
		<path d="M10 10h4"/>
	</svg>`

	LucideSimSecondary = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#64748b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
		<path d="M6 2h8l6 6v12a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2z"/>
		<path d="M10 10v4"/>
		<path d="M14 10v4"/>
		<path d="M10 14h4"/>
		<path d="M10 10h4"/>
	</svg>`

	LucideSmartphone = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#f45a51" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
		<rect width="14" height="20" x="5" y="2" rx="2" ry="2"/>
		<path d="M12 18h.01"/>
	</svg>`

	LucideShieldCheck = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#f45a51" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
		<path d="M20 13c0 5-3.5 7.5-7.66 8.95a1 1 0 0 1-.67-.01C7.5 20.5 4 18 4 13V6a1 1 0 0 1 1-1c2 0 4.5-1.2 6.24-2.72a1.17 1.17 0 0 1 1.52 0C14.51 3.8 17 5 19 5a1 1 0 0 1 1 1z"/>
		<path d="m9 12 2 2 4-4"/>
	</svg>`

	LucideUserCheck = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#f45a51" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
		<path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/>
		<circle cx="9" cy="7" r="4"/>
		<polyline points="16 11 18 13 22 9"/>
	</svg>`

	LucideFileCheck = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#10b981" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
		<path d="M15 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7Z"/>
		<path d="M14 2v4a2 2 0 0 0 2 2h4"/>
		<path d="m9 15 2 2 4-4"/>
	</svg>`

	LucideKeyRound = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#f45a51" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
		<path d="M2.586 17.414A2 2 0 0 0 2 18.828V21a1 1 0 0 0 1 1h3a1 1 0 0 0 1-1v-1a1 1 0 0 1 1-1h1a1 1 0 0 0 1-1v-1a1 1 0 0 1 1-1h.172a2 2 0 0 0 1.414-.586l.814-.814a6.5 6.5 0 1 0-4-4z"/>
		<circle cx="16.5" cy="7.5" r=".5" fill="#f45a51"/>
	</svg>`

	LucideLock = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#f45a51" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
		<rect width="18" height="11" x="3" y="11" rx="2" ry="2"/>
		<path d="M7 11V7a5 5 0 0 1 10 0v4"/>
	</svg>`

	LucideClock = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#f59e0b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
		<circle cx="12" cy="12" r="10"/>
		<polyline points="12 6 12 12 16 14"/>
	</svg>`

	LucideFingerprint = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#f45a51" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
		<path d="M12 10a2 2 0 0 0-2 2c0 1.02-.1 2.51-.26 4"/>
		<path d="M14 13.12c0 2.38 0 6.38-1 8.88"/>
		<path d="M17.29 21.02c.12-.6.43-2.3.5-3.02"/>
		<path d="M2 12a10 10 0 0 1 18-6"/>
		<path d="M2 16h.01"/>
		<path d="M21.8 16c.2-2 .131-5.354 0-6"/>
		<path d="M5 19.5C5.5 18 6 15 6 12a6 6 0 0 1 .34-2"/>
		<path d="M8.65 22c.21-.66.45-1.32.57-2"/>
		<path d="M9 6.8a6 6 0 0 1 9 5.2v2"/>
	</svg>`

	LucideIdCard = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#f45a51" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
		<path d="M16 10h2"/>
		<path d="M16 14h2"/>
		<path d="M6.17 15a3 3 0 0 1 5.66 0"/>
		<circle cx="9" cy="11" r="2"/>
		<rect x="2" y="5" width="20" height="14" rx="2"/>
	</svg>`

	LucideBadgeCheck = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#10b981" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
		<path d="M3.85 8.62a4 4 0 0 1 4.78-4.77 4 4 0 0 1 6.74 0 4 4 0 0 1 4.78 4.78 4 4 0 0 1 0 6.74 4 4 0 0 1-4.77 4.78 4 4 0 0 1-6.75 0 4 4 0 0 1-4.78-4.77 4 4 0 0 1 0-6.76Z"/>
		<path d="m9 12 2 2 4-4"/>
	</svg>`

	LucideAlertTriangle = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#f59e0b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
		<path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"/>
		<line x1="12" y1="9" x2="12" y2="13"/>
		<line x1="12" y1="17" x2="12.01" y2="17"/>
	</svg>`

	LucideAlertCircle = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#f45a51" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
		<circle cx="12" cy="12" r="10"/>
		<line x1="12" y1="8" x2="12" y2="12"/>
		<line x1="12" y1="16" x2="12.01" y2="16"/>
	</svg>`

	LucideCheckCircle2 = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#10b981" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
		<circle cx="12" cy="12" r="10"/>
		<path d="m9 12 2 2 4-4"/>
	</svg>`

	LucideInfo = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#38bdf8" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
		<circle cx="12" cy="12" r="10"/>
		<path d="M12 16v-4"/>
		<path d="M12 8h.01"/>
	</svg>`

	LucideFlag = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#f59e0b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
		<path d="M4 15s1-1 4-1 5 2 8 2 4-1 4-1V3s-1 1-4 1-5-2-8-2-4 1-4 1z"/>
		<line x1="4" y1="22" x2="4" y2="15"/>
	</svg>`

	LucideChevronRight = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#f2f6f8" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
		<path d="m9 18 6-6-6-6"/>
	</svg>`

	LucideArrowRight = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#ffffff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
		<path d="M5 12h14"/>
		<path d="m12 5 7 7-7 7"/>
	</svg>`

	LucideSparkles = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#f45a51" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
		<path d="m12 3-1.912 5.813a2 2 0 0 1-1.275 1.275L3 12l5.813 1.912a2 2 0 0 1 1.275 1.275L12 21l1.912-5.813a2 2 0 0 1 1.275-1.275L21 12l-5.813-1.912a2 2 0 0 1-1.275-1.275L12 3Z"/>
		<path d="M5 3v4"/>
		<path d="M19 17v4"/>
		<path d="M3 5h4"/>
		<path d="M17 19h4"/>
	</svg>`

	// Backward compatibility aliases
	SVGShieldCheck      = LucideShieldCheck
	SVGQRCodeFrame      = LucideQrCode
	SVGSIMCardActive    = LucideSim
	SVGSIMCardSecondary = LucideSimSecondary
	SVGFingerprint      = LucideFingerprint
	SVGClockGrace       = LucideClock
	SVGAcademicCap      = LucideBadgeCheck
	SVGCheckVerified    = LucideCheckCircle2
	SVGCampusLogo       = SVGBBDITLogo
)

// Resource helpers
func ResourceFromSVG(name, svgContent string) fyne.Resource {
	return fyne.NewStaticResource(name, []byte(svgContent))
}

// WhiteResourceFromSVG produces a pure white version of an SVG resource designed to complement primary action buttons.
func WhiteResourceFromSVG(name, svgContent string) fyne.Resource {
	re := strings.NewReplacer(
		`stroke="#f45a51"`, `stroke="#ffffff"`,
		`stroke="#10b981"`, `stroke="#ffffff"`,
		`stroke="#f59e0b"`, `stroke="#ffffff"`,
		`stroke="#64748b"`, `stroke="#ffffff"`,
		`stroke="#38bdf8"`, `stroke="#ffffff"`,
	)
	return fyne.NewStaticResource(name, []byte(re.Replace(svgContent)))
}

// RenderSVGImage creates an SVG canvas image with explicit bounding dimensions.
func RenderSVGImage(res fyne.Resource, width, height float32) *canvas.Image {
	img := canvas.NewImageFromResource(res)
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(width, height))
	return img
}
