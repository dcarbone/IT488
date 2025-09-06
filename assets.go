package main

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/fs"
	"path"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"golang.org/x/image/draw"
)

var (
	//go:embed assets
	appAssets embed.FS
)

var (
	AssetImageLogo    = MustGetAssetImage("logo.png")
	AssetImageWarning = MustGetAssetImage("warning.png")

	AssetImagePriorityLowest  = MustGetAssetImage("lowest_priority.png")
	AssetImagePriorityLow     = MustGetAssetImage("low_priority.png")
	AssetImagePriorityHigh    = MustGetAssetImage("high_priority.png")
	AssetImagePriorityHighest = MustGetAssetImage("highest_priority.png")

	AssetImageStatusTodo = MustGetAssetImage("status_todo.png")
	AssetImageStatusDone = MustGetAssetImage("status_done.png")
	AssetImageStatusSkip = MustGetAssetImage("status_skip.png")
)

func GetAssetFile(name string) (fs.File, error) {
	return appAssets.Open(path.Join("assets", name))
}

func GetAssetBytes(name string) ([]byte, error) {
	f, err := GetAssetFile(name)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	return io.ReadAll(f)
}

func GetAssetImage(name string) (image.Image, error) {
	f, err := GetAssetFile(name)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	img, _, err := image.Decode(f)
	return img, err
}

func MustGetAssetImage(name string) image.Image {
	img, err := GetAssetImage(name)
	if err != nil {
		panic(fmt.Sprintf("Error retreiving image %q: %v", name, err))
	}
	return img
}

func ResizePNG(src image.Image, scale float64) image.Image {
	dst := image.NewRGBA(
		image.Rect(0, 0, int(float64(src.Bounds().Max.X)/scale), int(float64(src.Bounds().Max.Y)/scale)),
	)
	draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)
	return dst
}

func GetConstrainedImage(img image.Image, maxDimension float64) image.Image {
	var scale float64
	bounds := img.Bounds()
	if bounds.Max.X > bounds.Max.Y {
		scale = float64(bounds.Max.X) / maxDimension
	} else {
		scale = float64(bounds.Max.Y) / maxDimension
	}
	return ResizePNG(img, scale)
}

func GetAssetImageCanvas(src image.Image, opts ...func(mg *canvas.Image)) *canvas.Image {
	img := canvas.NewImageFromImage(src)
	img.FillMode = canvas.ImageFillOriginal
	for _, opt := range opts {
		opt(img)
	}
	return img
}

func EncodeImageToResource(name string, img image.Image) *fyne.StaticResource {
	buf := bytes.NewBuffer(nil)
	err := png.Encode(buf, img)
	if err != nil {
		panic(fmt.Sprintf("error encoding %q image: %v", name, err))
	}
	return fyne.NewStaticResource(name, buf.Bytes())
}
