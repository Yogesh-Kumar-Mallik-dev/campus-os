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

// ShadcnButton extends Fyne's Button with shadcn variant tokens, accessible hit targets,
// and full interactive states (Hover, Active, Focus, and Disabled).
type ShadcnButton struct {
	widget.Button
	Variant    ButtonVariant
	ButtonSize ButtonSize
}

// NewShadcnButton constructs a button adhering to shadcn variants and sizing.
func NewShadcnButton(label string, variant ButtonVariant, size ButtonSize, icon fyne.Resource, onClick func()) *ShadcnButton {
	btn := &ShadcnButton{
		Variant:    variant,
		ButtonSize: size,
	}
	btn.Text = label
	btn.Icon = icon
	btn.OnTapped = onClick

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

	btn.ExtendBaseWidget(btn)
	return btn
}

// MinSize enforces Guardrail 26: comfortable minimum touch and click targets (>= 36-38px height).
func (b *ShadcnButton) MinSize() fyne.Size {
	base := b.Button.MinSize()
	targetHeight := float32(38)
	targetWidth := base.Width

	switch b.ButtonSize {
	case ButtonSizeSm:
		targetHeight = 32
	case ButtonSizeLg:
		targetHeight = 44
	case ButtonSizeIcon:
		targetHeight = 38
		if targetWidth < 38 {
			targetWidth = 38
		}
	default:
		targetHeight = 38
	}

	if base.Height < targetHeight {
		base.Height = targetHeight
	}
	if base.Width < targetWidth {
		base.Width = targetWidth
	}
	return base
}

