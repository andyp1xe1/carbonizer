package main

import (
	"flag"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	"carbonizer/filters"
	q "carbonizer/filters/quantizer"
	"carbonizer/utils"
)

var pathFlag *string
var quantizeDepthFlag *int

func main() {
	pathFlag = flag.String("file", "-", "the input file")
	quantizeDepthFlag = flag.Int(
		"quantizer", 0, "the detph of the median cut quantization")

	flag.Parse()

	f, err := os.Open(*pathFlag)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	img, err := decodeImage(f)
	if err != nil {
		log.Fatal(err)
	}

	layers := filtersFromFlags()

	canvas := filters.NewCanvas(img, layers...)
	res, err := canvas.Calc()

	if err != nil {
		log.Fatal(err)
	}

	fbase := filepath.Base(f.Name())
	fname := strings.Split(fbase, ".")[0]

	resName, err := utils.SaveImage(fname, res)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("result path:", resName)
	}

}

func filtersFromFlags() []filters.Filter {
	var filters []filters.Filter

	if *quantizeDepthFlag != 0 {
		quantizer := q.NewFilter(
			*quantizeDepthFlag,
			//q.WithDither(q.Bayer(3)),
			q.WithDither(q.FloydSteinsberg),
			//q.WithStdPalette,
			//q.WithMapPalette,
			//q.WithMemoPalette,
		)
		filters = append(filters, quantizer)
	}
	// More filters here

	return filters
}

func decodeImage(f *os.File) (image.Image, error) {
	fext := filepath.Ext(filepath.Ext(f.Name()))

	var img image.Image
	var err error

	switch fext {
	case "jpeg":
		img, err = jpeg.Decode(f)
	case "png":
		img, err = png.Decode(f)
	default:
		img, _, err = image.Decode(f)
	}

	if err != nil {
		return nil, err
	}

	return img, err
}
