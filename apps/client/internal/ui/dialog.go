package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
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
		cancelBtn,
		confirmBtn,
	)

	card := NewShadcnCard(CardParts{
		Title:       title,
		Description: description,
		Footer:      container.NewBorder(nil, nil, nil, buttons),
	})

	popup = widget.NewModalPopUp(container.NewGridWrap(fyne.NewSize(420, card.MinSize().Height), card), w.Canvas())
	popup.Show()
	return popup
}
