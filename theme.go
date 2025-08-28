package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	ThemeBackgroundColorHex = "#f2df79"
	ThemeBackgroundColorRGB = "rgb(242, 223, 121)"
)

var (
	ColorRed = color.RGBA{R: 255}
)

func ThemeBackgroundColor() color.Color {
	return color.RGBA{
		R: 242,
		G: 223,
		B: 121,
	}
}

var _ fyne.Theme = (*TodoTodayTheme)(nil)

type TodoTodayTheme struct {
	fyne.Theme
}

func NewTheme() *TodoTodayTheme {
	th := TodoTodayTheme{
		Theme: theme.DefaultTheme(),
	}
	return &th
}

func (th *TodoTodayTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return ThemeBackgroundColor()

	default:
		return th.Theme.Color(name, theme.VariantLight)
	}
}

func HeaderCanvas(text string, opts ...func(txt *canvas.Text)) *canvas.Text {
	txt := canvas.NewText(text, color.Black)
	txt.Alignment = fyne.TextAlignCenter
	txt.TextSize = 32
	txt.TextStyle = fyne.TextStyle{
		Bold: true,
	}
	for _, opt := range opts {
		opt(txt)
	}
	return txt
}

func FormLabel(text string, opts ...func(lbl *widget.Label)) *widget.Label {
	lbl := widget.NewLabel(text)
	lbl.TextStyle = fyne.TextStyle{
		Bold: true,
	}
	for _, opt := range opts {
		opt(lbl)
	}
	return lbl
}
