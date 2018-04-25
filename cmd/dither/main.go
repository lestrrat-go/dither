package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	dither "github.com/lestrrat-go/dither"
	"github.com/pkg/errors"
)

func main() {
	if err := _main(); err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}
}

type file struct {
	name string
}

// Command line flags
var (
	filterList string
	outputDir  string
	export     string
	grayscale  bool
	threshold  bool
	multiplier float64
	commands   flag.FlagSet
)

const helper = `
Usage:
go run <image>

Options:
  -export string
    	Generate the color and greyscale dithered images. Options: 'all', 'color', 'mono' (default "all")
  -grayscale
    	Convert image to grayscale (default true)
  -multiplier float
    	Error multiplier (default 1.18)
  -outputdir string
    	Directory name, where to save the generated images (default "output")
  -threshold
    	Export threshold image (default true)
`

// Open the input file
func (file *file) Open() (image.Image, error) {
	f, err := os.Open(file.name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	return img, errors.Wrapf(err, `failed to decode image '%s'`, file.name)
}

func writePNG(im image.Image, fn string) error {
	f, err := os.Create(fn)
	if err != nil {
		return errors.Wrapf(err, `failed to create file: %s`, fn)
	}
	defer f.Close()

	if err := png.Encode(f, im); err != nil {
		return errors.Wrapf(err, `failed to encode file %s`, fn)
	}
	return nil
}

func _main() error {
	commands = *flag.NewFlagSet("commands", flag.ExitOnError)

	commands.StringVar(&filterList, "filters", "all", "Comma-separated names of filters to apply. Use 'all' to apply all filters")
	commands.StringVar(&outputDir, "outputdir", "output", "Directory name, where to save the generated images")
	commands.StringVar(&export, "export", "all", "Generate the color and greyscale dithered images. Options: 'all', 'color', 'mono'")
	commands.BoolVar(&grayscale, "grayscale", true, "Convert image to grayscale")
	commands.BoolVar(&threshold, "threshold", true, "Export threshold image")
	commands.Float64Var(&multiplier, "multiplier", 1.18, "Error multiplier")

	if len(os.Args) <= 1 || (os.Args[1] == "--help" || os.Args[1] == "-h") {
		fmt.Println(errors.New(helper))
		os.Exit(1)
	}

	// Parse flags before to use them
	commands.Parse(os.Args[2:])

	var filters []*dither.Filter
	for _, fname := range strings.Split(filterList, ",") {
		switch strings.TrimSpace(fname) {
		case "all":
			filters = []*dither.Filter{
				dither.Atkinson,
				dither.Burkes,
				dither.FloydSteinberg,
				dither.Stucki,
				dither.Sierra2,
				dither.Sierra3,
				dither.SierraLite,
			}
			break
		case "atkinson":
			filters = append(filters, dither.Atkinson)
		case "burkes":
			filters = append(filters, dither.Burkes)
		case "floyd-steinberg":
			filters = append(filters, dither.FloydSteinberg)
		case "stucki":
			filters = append(filters, dither.Stucki)
		case "sierra2":
			filters = append(filters, dither.Sierra2)
		case "sierra3":
			filters = append(filters, dither.Sierra3)
		case "sierra-lite":
			filters = append(filters, dither.SierraLite)
		}
	}

	// Channel to signal the completion event
	src, err := os.Open(os.Args[1])
	if err != nil {
		return errors.Wrapf(err, `failed to open source file %s`, os.Args[1])
	}
	defer src.Close()

	img, _, err := image.Decode(src)
	if err != nil {
		return errors.Wrapf(err, `failed to decode image '%s'`, os.Args[1])
	}

	fmt.Print("Applying filters [")
	for i, filter := range filters {
		fmt.Print(filter.Name())
		if i < len(filters)-1 {
			fmt.Print(", ")
		}
	}
	fmt.Println("]")

	fmt.Print("Rendering image...")
	now := time.Now()

	done := make(chan struct{})

	exportDirs := []string{
		filepath.Join(outputDir, "mono"),
		filepath.Join(outputDir, "color"),
	}

	for _, subdir := range exportDirs {
		if err := os.MkdirAll(subdir, os.ModePerm); err != nil {
			return errors.Wrapf(err, `failed to create output directory %s`, subdir)
		}
	}

	gray := dither.Grayscale(img)
	if grayscale {
		if err := writePNG(gray, filepath.Join(outputDir, "grayscale.png")); err != nil {
			return errors.Wrap(err, `failed to write grayscale file`)
		}
	}

	if threshold {
		if err := writePNG(dither.Threshold(gray), filepath.Join(outputDir, "threshold.png")); err != nil {
			return errors.Wrap(err, `failed to write threshold file`)
		}
	}

	// Function to visualize the rendering progress
	go func(done chan struct{}) {
		ticker := time.NewTicker(time.Millisecond * 200)
		for {
			select {
			case <-ticker.C:
				fmt.Print(".")
			case <-done:
				ticker.Stop()
			}
		}
	}(done)

	var exportColor, exportMono bool
	switch export {
	case "all":
		exportColor = true
		exportMono = true
	case "color":
		exportColor = true
	case "mono":
		exportMono = true
	}

	// Run dither methods
	var wg sync.WaitGroup
	for _, filter := range filters {
		wg.Add(1)
		go func(filter *dither.Filter) {
			defer wg.Done()

			if exportMono {
				outputMono := dither.Monochrome(filter.Matrix(), img, float32(multiplier))
				writePNG(outputMono, filepath.Join(exportDirs[0], filter.Name()+".png"))
			}
			if exportColor {
				outputColor := dither.Color(filter.Matrix(), img, float32(multiplier))
				writePNG(outputColor, filepath.Join(exportDirs[1], filter.Name()+".png"))
			}
		}(filter)
	}
	wg.Wait()
	close(done)

	since := time.Since(now)
	fmt.Println("\nDone✓")
	fmt.Printf("Rendered in: %.2fs\n", since.Seconds())
	return nil
}

// Output the resulting image
func generateOutput(f *dither.Filter, img image.Image, exportDir string) {
	output, err := os.Create(filepath.Join(exportDir, f.Name()+".png"))
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()

	err = png.Encode(output, img)
	if err != nil {
		log.Fatal(err)
	}
}
