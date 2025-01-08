package quantizer

import (
	"image"
	"image/color"
)

type colorVec [3]uint8

func (c colorVec) R() uint8 {
	return c[0]
}

func (c colorVec) G() uint8 {
	return c[1]
}

func (c colorVec) B() uint8 {
	return c[2]
}

func span(b []colorVec) int {
	var span [3]uint8

	for axis := 0; axis < 3; axis++ {
		minLum := b[0][axis]
		maxLum := minLum

		bucket := b[1:]
		for _, col := range bucket {
			lum := col[axis]

			minLum = min(minLum, lum)
			maxLum = max(maxLum, lum)
		}

		span[axis] = maxLum - minLum
	}

	maxSpan := max(span[0], span[1], span[2])
	var i int
	for i = 0; i < 3; i++ {
		if span[i] == maxSpan {
			break
		}
	}
	return i
}

func colorVecAt(src image.Image, x, y int) colorVec {
	var vec colorVec
	var tcol color.Color

	switch i := src.(type) {
	case *image.YCbCr:
		yi := i.YOffset(x, y)
		ci := i.COffset(x, y)
		c := color.YCbCr{
			i.Y[yi],
			i.Cb[ci],
			i.Cr[ci],
		}
		tcol = c
		vec = colorVec{c.Y, c.Cb, c.Cr}
	case *image.RGBA:
		ci := i.PixOffset(x, y)
		tcol = color.RGBA{i.Pix[ci+0], i.Pix[ci+1], i.Pix[ci+2], 255}
		vec = colorVec{i.Pix[ci+0], i.Pix[ci+1], i.Pix[ci+2]}
	default:
		tcol = i.At(x, y)
		col := color.RGBAModel.Convert(tcol).(color.RGBA)
		vec = colorVec{col.R, col.G, col.B}
	}

	return vec
}
