package ui

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// ShowToast displays a floating transient notification overlay mirroring shadcn sonner.
func ShowToast(w fyne.Window, title, message string, variant AlertVariant, duration time.Duration) *widget.PopUp {
	if duration <= 0 {
		duration = 3 * time.Second
	}

	var icon fyne.Resource
	switch variant {
	case AlertSuccess:
		icon = ResourceFromSVG("toast_check.svg", SVGCheckVerified)
	case AlertWarning:
		icon = ResourceFromSVG("toast_warn.svg", SVGClockGrace)
	case AlertDestructive:
		icon = ResourceFromSVG("toast_err.svg", SVGClockGrace)
	default:
		icon = ResourceFromSVG("toast_info.svg", SVGShieldCheck)
	}

	alertCard := NewShadcnAlert(title, message, variant, icon)

	popup := widget.NewPopUp(alertCard, w.Canvas())
	popup.Move(fyne.NewPos((w.Canvas().Size().Width-alertCard.MinSize().Width)/2, 40))

	time.AfterFunc(duration, func() {
		popup.Hide()
	})

	return popup
}
