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

func Color(m Matrixer, input image.Image, errorMultiplier float32) image.Image {
	bounds := input.Bounds()
	img := image.NewRGBA(bounds)
	for x := bounds.Min.X; x < bounds.Dx(); x++ {
		for y := bounds.Min.Y; y < bounds.Dy(); y++ {
			pixel := input.At(x, y)
			img.Set(x, y, pixel)
		}
	}
	dx, dy := bounds.Dx(), bounds.Dy()

	// Prepopulate multidimensional slices
	redErrors := make([][]float32, dx)
	greenErrors := make([][]float32, dx)
	blueErrors := make([][]float32, dx)
	for x := 0; x < dx; x++ {
		redErrors[x] = make([]float32, dy)
		greenErrors[x] = make([]float32, dy)
		blueErrors[x] = make([]float32, dy)
		for y := 0; y < dy; y++ {
			redErrors[x][y] = 0
			greenErrors[x][y] = 0
			blueErrors[x][y] = 0
		}
	}

	// Diffuse error in two dimension
	matrix := m.Matrix()
	ydim := matrix.Rows() - 1
	xdim := matrix.Cols() / 2
	var qrr, qrg, qrb float32
	for x := 0; x < dx; x++ {
		for y := 0; y < dy; y++ {
			r32, g32, b32, a := img.At(x, y).RGBA()
			r, g, b := float32(uint8(r32)), float32(uint8(g32)), float32(uint8(b32))
			r -= redErrors[x][y] * errorMultiplier
			g -= greenErrors[x][y] * errorMultiplier
			b -= blueErrors[x][y] * errorMultiplier

			// Diffuse the error of each calculation to the neighboring pixels
			if r < 128 {
				qrr = -r
				r = 0
			} else {
				qrr = 255 - r
				r = 255
			}
			if g < 128 {
				qrg = -g
				g = 0
			} else {
				qrg = 255 - g
				g = 255
			}
			if b < 128 {
				qrb = -b
				b = 0
			} else {
				qrb = 255 - b
				b = 255
			}
			img.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})

			for xx := 0; xx < ydim+1; xx++ {
				for yy := -xdim; yy <= xdim-1; yy++ {
					if y+yy < 0 || dy <= y+yy || x+xx < 0 || dx <= x+xx {
						continue
					}
					// Adds the error of the previous pixel to the current pixel
					factor := matrix.Get(yy+ydim, xx)
					redErrors[x+xx][y+yy] += qrr * factor
					greenErrors[x+xx][y+yy] += qrg * factor
					blueErrors[x+xx][y+yy] += qrb * factor
				}
			}
		}
	}
	return img
}

func Monochrome(m Matrixer, input image.Image, errorMultiplier float32) image.Image {
	bounds := input.Bounds()
	img := image.NewGray(bounds)
	for x := bounds.Min.X; x < bounds.Dx(); x++ {
		for y := bounds.Min.Y; y < bounds.Dy(); y++ {
			pixel := input.At(x, y)
			img.Set(x, y, pixel)
		}
	}
	dx, dy := bounds.Dx(), bounds.Dy()

	// Prepopulate multidimensional slice
	errors := NewMatrix(dx, dy)

	matrix := m.Matrix()
	ydim := matrix.Rows() - 1
	xdim := matrix.Cols() / 2
	for x := 0; x < dx; x++ {
		for y := 0; y < dy; y++ {
			pix := float32(img.GrayAt(x, y).Y)
			pix -= errors.Get(x, y) * errorMultiplier

			var quantError float32
			// Diffuse the error of each calculation to the neighboring pixels
			if pix < 128 {
				quantError = -pix
				pix = 0
			} else {
				quantError = 255 - pix
				pix = 255
			}

			img.SetGray(x, y, color.Gray{Y: uint8(pix)})

			// Diffuse error in two dimension
			for xx := 0; xx < ydim+1; xx++ {
				for yy := -xdim; yy <= xdim-1; yy++ {
					if y+yy < 0 || dy <= y+yy || x+xx < 0 || dx <= x+xx {
						continue
					}
					// Adds the error of the previous pixel to the current pixel
					prev := errors.Get(x+xx, y+yy)
					delta := quantError * matrix.Get(yy+ydim, xx)
					errors.Set(x+xx, y+yy, prev+delta)
				}
			}
		}
	}
	return img
}
