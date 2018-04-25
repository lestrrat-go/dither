package dither

import (
	"image"
	"image/color"
)

var white = color.Gray{Y: 255}
var black = color.Gray{Y: 0}

func threshold(pixel color.Gray) color.Gray {
	if pixel.Y > 123 {
		return white
	}
	return black
}

func Threshold(input *image.Gray) *image.Gray {
	bounds := input.Bounds()
	dithered := image.NewGray(bounds)
	dx := bounds.Dx()
	dy := bounds.Dy()
	for x := 0; x < dx; x++ {
		for y := 0; y < dy; y++ {
			dithered.Set(x, y, threshold(input.GrayAt(x, y)))
		}
	}
	return dithered
}

func Grayscale(input image.Image) *image.Gray {
	bounds := input.Bounds()
	gray := image.NewGray(bounds)
	dx := bounds.Dx()
	dy := bounds.Dy()

	for x := bounds.Min.X; x < dx; x++ {
		for y := bounds.Min.Y; y < dy; y++ {
			gray.Set(x, y, input.At(x, y))
		}
	}
	return gray
}
