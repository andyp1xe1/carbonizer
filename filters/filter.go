package filters

import (
	"fmt"
	"image"
	"image/draw"
)

type Filter interface {
	Transform(in *image.RGBA) error
}

type Canvas struct {
	src    image.Image
	Layers []Filter
}

func NewCanvas(img image.Image, layers ...Filter) *Canvas {
	return &Canvas{img, layers}
}

func (c Canvas) Calc() (image.Image, error) {
	var err error

	in := image.NewRGBA(c.src.Bounds())
	draw.Draw(in, c.src.Bounds(), c.src, c.src.Bounds().Min, draw.Src)

	for i, filter := range c.Layers {
		if err = filter.Transform(in); err != nil {
			return in, fmt.Errorf("failed and stopped at layer (%v): %w", i, err)
		}
	}

	return in, nil
}
