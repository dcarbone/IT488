package main

import (
	"embed"
	"image"
	"io"
	"io/fs"
	"path"

	"golang.org/x/image/draw"
)

var (
	//go:embed assets
	appAssets embed.FS
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

func GetFullSizeLogoPNG() (image.Image, error) {
	return GetAssetImage("todo_today_transparent_logo.png")
}

func GetConstrainedLogoPNG() (image.Image, error) {
	src, err := GetFullSizeLogoPNG()
	if err != nil {
		return nil, err
	}

	scale := src.Bounds().Max.X / 400

	dst := image.NewRGBA(image.Rect(0, 0, src.Bounds().Max.X/scale, src.Bounds().Max.Y/scale))
	draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)
	return dst, nil
}

func ResizePNG(src image.Image, scale int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, src.Bounds().Max.X/scale, src.Bounds().Max.Y/scale))
	draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)
	return dst
}
