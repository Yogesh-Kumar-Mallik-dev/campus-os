package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// ButtonVariant enumerates shadcn button style variants.
type ButtonVariant int

const (
	ButtonDefault ButtonVariant = iota
	ButtonSecondary
	ButtonDestructive
	ButtonOutline
	ButtonGhost
)

// ButtonSize controls button vertical padding and font size.
type ButtonSize int

const (
	ButtonSizeDefault ButtonSize = iota
	ButtonSizeSm
	ButtonSizeLg
	ButtonSizeIcon
)

// NewShadcnButton constructs a button adhering to shadcn variants and sizing.
func NewShadcnButton(label string, variant ButtonVariant, size ButtonSize, icon fyne.Resource, onClick func()) *widget.Button {
	var btn *widget.Button

	if icon != nil && label != "" {
		btn = widget.NewButtonWithIcon(label, icon, onClick)
	} else if icon != nil {
		btn = widget.NewButtonWithIcon("", icon, onClick)
	} else {
		btn = widget.NewButton(label, onClick)
	}

	switch variant {
	case ButtonDefault:
		btn.Importance = widget.HighImportance
	case ButtonSecondary:
		btn.Importance = widget.MediumImportance
	case ButtonDestructive:
		btn.Importance = widget.DangerImportance
	case ButtonOutline, ButtonGhost:
		btn.Importance = widget.LowImportance
	}

	return btn
}
