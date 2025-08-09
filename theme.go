package main

import (
	"image/color"
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
