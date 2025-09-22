package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

var (
	ColorRed = color.RGBA{R: 255}

	ColorBackground = color.RGBA{R: 242, G: 223, B: 121} // F2DF79
	ColorBlue       = color.RGBA{R: 11, G: 2, B: 133}    // 0B0285
	ColorPurple     = color.RGBA{R: 139, G: 129, B: 253} // 8B81FD
	ColorPink       = color.RGBA{R: 227, G: 204, B: 252} // E3CCFC
	ColorYellow     = color.RGBA{R: 221, G: 253, B: 204} // DDFDCC
	ColorGreen      = color.RGBA{R: 2, G: 116, B: 102}   // 027466

	IconEdit = EncodeImageToResource("edit", AssetImageEditIcon)
)

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
		return ColorBackground

	default:
		return th.Theme.Color(name, theme.VariantLight)
	}
}

func ResizeTextToFit(txt *canvas.Text, baseSize, maxWidth float32) {
	txt.TextSize = baseSize
	actual := fyne.MeasureText(txt.Text, txt.TextSize, txt.TextStyle)
	for actual.Width > maxWidth && txt.TextSize > 8 {
		txt.TextSize -= 1
		actual = fyne.MeasureText(txt.Text, txt.TextSize, txt.TextStyle)
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

func FormLabel(text string, opts ...func(txt *canvas.Text)) *canvas.Text {
	txt := canvas.NewText(text, color.Black)
	txt.TextStyle = fyne.TextStyle{
		Bold: true,
	}
	ResizeTextToFit(txt, 14, 150)
	for _, opt := range opts {
		opt(txt)
	}
	return txt
}
