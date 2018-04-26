package dither_test

import (
	"image"
	"image/png"
	"os"

	"github.com/lestrrat-go/dither"
)

func Example() {
	f, err := os.Open("file.png")
	if err != nil {
		return
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return
	}

	ditheredImg := dither.Monochrome(dither.Burkes, img, 1.18)

	png.Encode(os.Stdout, ditheredImg)
}
