package quantizer

import (
	"image"
)

type Filter struct {
	*quantizer
}

func NewFilter(depth int, opt paletteOpt) *Filter {
	return &Filter{
		newQuantizer(depth, opt),
	}
}

func (f *Filter) Transform(in *image.RGBA) error {
	f.medianCut(in)

	return f.quantizeImg(in)
}
