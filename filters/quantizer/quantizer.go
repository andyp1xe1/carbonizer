package quantizer

import (
	"image"
)

type paletteMaker func(q *quantizer) Palette

type quantizerFunc func(q *quantizer, img *image.RGBA)

var (
	MedianCut quantizerFunc = medianCut
)

type paletteOpt struct {
	makePalette paletteMaker
	ditherer    DitherIter
}

type DitherIter func(img *image.RGBA, p Palette, x, y int)

var FloydSteinsberg DitherIter = floydSteinsbergIter
var Bayer = makeBayerIter

func WithDither(iter DitherIter) paletteOpt {
	return paletteOpt{
		paletteMaker(withMemoPalette),
		iter,
	}
}

var (
	WithMemoPalette = paletteOpt{withMemoPalette, noDither}
	WithMapPalette  = paletteOpt{withMapPalette, noDither}
	WithStdPalette  = paletteOpt{withStdPalette, noDither}
)

type quantizer struct {
	colorBucket
	depth int

	quantize   quantizerFunc
	paletteOpt paletteOpt
}

func newQuantizer(depth int, qf quantizerFunc, opt paletteOpt) *quantizer {
	return &quantizer{
		colorBucket: newBucket(),
		depth:       depth,
		quantize:    qf,
		paletteOpt:  opt,
	}
}

func medianCut(q *quantizer, img *image.RGBA) {
	q.makeBucketBuff(img)
	currBucket := colorBucket{q.buff, q.axis}

	depth := uint(q.depth)
	for lvl := uint(0); lvl < depth; lvl++ {
		for currIdx := uint(0); currIdx < (1 << lvl); currIdx++ {
			start := q.getMarker(currIdx, lvl)
			end := q.getMarker(currIdx+1, lvl)

			currBucket.buff = q.buff[start:end]
			currBucket.updateAxis()
			currBucket.sort()
		}
	}
}

func (q *quantizer) exportImg(img *image.RGBA) error {
	bounds := img.Bounds()

	p := q.paletteOpt.makePalette(q)
	dither := q.paletteOpt.ditherer

	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			dither(img, p, x, y)
		}
	}
	return nil
}

func noDither(img *image.RGBA, p Palette, x, y int) {
	col := colorVecAt(img, x, y)
	newCol := p.convVec(col)
	img.SetRGBA(x, y, newCol)
}

func withMapPalette(q *quantizer) Palette {
	palette := make(MapPalette)

	depth := uint(q.depth)
	for currIdx := uint(0); currIdx < (1 << q.depth); currIdx++ {
		start := q.getMarker(currIdx, depth)
		end := q.getMarker(currIdx+1, depth)

		avgCol := q.avgPart(start, end)

		for _, col := range q.buff[start:end] {
			palette[col] = avgCol
		}

	}

	return palette
}

func withStdPalette(q *quantizer) Palette {
	return calcStdPalette(q)
}

func withMemoPalette(q *quantizer) Palette {
	return newMemo(calcStdPalette(q))
}

func calcStdPalette(q *quantizer) StdPalette {
	var palette StdPalette

	depth := uint(q.depth)
	for currIdx := uint(0); currIdx < (1 << q.depth); currIdx++ {
		start := q.getMarker(currIdx, depth)
		end := q.getMarker(currIdx+1, depth)

		avgCol := q.avgPart(start, end)

		palette = append(palette, avgCol)

	}

	return palette
}
