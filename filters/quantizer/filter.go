package quantizer

import (
	"image"
)

type Filter struct {
	*quantizer
}

func NewFilter(depth int, qf quantizerFunc, paletteOpt paletteOpt) *Filter {
	return &Filter{
		newQuantizer(depth, qf, paletteOpt),
	}
}

func (f *Filter) Transform(in *image.RGBA) error {
	f.quantize(f.quantizer, in)

	return f.exportImg(in)
}
