package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// ShowAlertDialog renders a modal confirmation dialog overlay matching shadcn alert-dialog.
func ShowAlertDialog(w fyne.Window, title, description, confirmText string, isDestructive bool, onConfirm func()) *widget.PopUp {
	var popup *widget.PopUp

	cancelBtn := NewShadcnButton("Cancel", ButtonOutline, ButtonSizeDefault, nil, func() {
		if popup != nil {
			popup.Hide()
		}
	})

	btnVariant := ButtonDefault
	if isDestructive {
		btnVariant = ButtonDestructive
	}

	confirmBtn := NewShadcnButton(confirmText, btnVariant, ButtonSizeDefault, nil, func() {
		if popup != nil {
			popup.Hide()
		}
		if onConfirm != nil {
			onConfirm()
		}
	})

	buttons := container.NewHBox(
		layout.NewSpacer(),
		cancelBtn,
		confirmBtn,
	)

	card := NewShadcnCard(CardParts{
		Title:       title,
		Description: description,
		Footer:      buttons,
	})

	cardMin := card.MinSize()
	targetWidth := float32(440)
	targetHeight := cardMin.Height
	if targetHeight < 200 {
		targetHeight = 200
	}

	popup = widget.NewModalPopUp(container.NewGridWrap(fyne.NewSize(targetWidth, targetHeight), card), w.Canvas())
	popup.Show()
	return popup
}
