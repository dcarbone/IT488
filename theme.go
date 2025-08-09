package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

const (
	ThemeBackgroundColorHex = "#f2df79"
	ThemeBackgroundColorRGB = "rgb(242, 223, 121)"
)

func ThemeBackgroundColor() color.Color {
	return color.RGBA{
		R: 242,
		G: 223,
		B: 121,
	}
}

var _ fyne.Theme = (*BackgroundColorTheme)(nil)

type BackgroundColorTheme struct {
	fyne.Theme
}

func NewTheme() *BackgroundColorTheme {
	th := BackgroundColorTheme{
		Theme: theme.DefaultTheme(),
	}
	return &th
}

func (th *BackgroundColorTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return ThemeBackgroundColor()

	default:
		return th.Theme.Color(name, variant)
	}
}
