package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// FormField encapsules a responsive labeled input field.
type FormField struct {
	Label       string
	Placeholder string
	HelperText  string
	IsPassword  bool
	Entry       *widget.Entry
}

// NewShadcnInput creates an entry widget configured with shadcn input styling.
func NewShadcnInput(placeholder string, isPassword bool) *widget.Entry {
	var entry *widget.Entry
	if isPassword {
		entry = widget.NewPasswordEntry()
	} else {
		entry = widget.NewEntry()
	}
	entry.SetPlaceHolder(placeholder)
	return entry
}

// NewFormField constructs a complete responsive field with label, input box, and helper text.
func NewFormField(field FormField) (fyne.CanvasObject, *widget.Entry) {
	entry := field.Entry
	if entry == nil {
		entry = NewShadcnInput(field.Placeholder, field.IsPassword)
	}

	items := make([]fyne.CanvasObject, 0, 3)

	if field.Label != "" {
		lbl := widget.NewLabelWithStyle(field.Label, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		lbl.Wrapping = fyne.TextWrapWord
		items = append(items, lbl)
	}

	items = append(items, entry)

	if field.HelperText != "" {
		helper := canvas.NewText(field.HelperText, color.NRGBA{R: 153, G: 167, B: 173, A: 220})
		helper.TextSize = 11.5
		items = append(items, helper)
	}

	return container.NewVBox(items...), entry
}
