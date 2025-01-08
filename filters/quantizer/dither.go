package quantizer

import (
	"image"
	"image/color"
)

type DitherMatrix struct {
	w, h, dfac int
	data       [][]int
}

func bayerUp(mat *DitherMatrix) DitherMatrix {
	fac := 4
	up := DitherMatrix{
		w: mat.w * 2, h: mat.h * 2,
		dfac: mat.dfac * fac,
	}
	up.data = make([][]int, up.h)
	for i := range up.h {
		up.data[i] = make([]int, up.w)
	}
	for y := range mat.h {
		for x := range mat.w {
			newCell := fac * mat.data[x][y]
			up.data[y][x] = newCell
			up.data[y][x+mat.w] = newCell + 2
			up.data[y+mat.h][x] = newCell + 3
			up.data[y+mat.h][x+mat.w] = newCell + 1
		}
	}
	return up
}

func makeBayerIter(n int) func(img *image.RGBA, p Palette, x, y int) {
	bayer := DitherMatrix{
		w: 2, h: 2, dfac: 4,
		data: [][]int{
			{0, 2},
			{3, 1},
		},
	}
	for range n {
		bayer = bayerUp(&bayer)
	}
	return func(img *image.RGBA, p Palette, x, y int) {
		xm, ym := x%bayer.w, y%bayer.h
		d := uint8(bayer.data[ym][xm] * 256 / bayer.dfac)

		col := colorVecAt(img, x, y)
		for i := range col {
			col[i] = clampLum(int(col[i]) + int(d) - 127)
		}

		img.SetRGBA(x, y, p.convVec(col))
	}

}

func floydSteinsbergIter(img *image.RGBA, p Palette, x, y int) {
	var (
		points = []image.Point{
			{x + 1, y}, {x + 1, y + 1}, {x, y + 1}, {x - 1, y + 1}}

		weights = []int{7, 1, 5, 3}
	)

	currCol := colorAt(img, image.Point{x, y})
	qtdCol := p.conv(currCol)

	errCol := diffColor(currCol, qtdCol)

	for i := range 4 {
		if !(points[i].In(img.Rect)) {
			continue
		}
		img.SetRGBA(
			points[i].X, points[i].Y,
			weightColor(colorAt(img, points[i]), errCol, weights[i]),
		)
	}

	img.SetRGBA(x, y, qtdCol)
}

func colorAt(img *image.RGBA, pt image.Point) color.RGBA {
	ci := img.PixOffset(pt.X, pt.Y)
	return color.RGBA{img.Pix[ci+0], img.Pix[ci+1], img.Pix[ci+2], 255}
}

type diffCol [3]int

func diffColor(c1, c2 color.RGBA) diffCol {
	return diffCol{
		int(c1.R) - int(c2.R),
		int(c1.G) - int(c2.G),
		int(c1.B) - int(c2.B),
	}
}

func weightColor(c1 color.RGBA, diff diffCol, factor int) color.RGBA {
	return color.RGBA{
		R: clampLum(int(c1.R) + (diff[0]*factor)/16),
		G: clampLum(int(c1.G) + (diff[1]*factor)/16),
		B: clampLum(int(c1.B) + (diff[2]*factor)/16),
		A: 255,
	}
}

func clampLum(val int) uint8 {
	if val < 0 {
		return 0
	}
	if val > 255 {
		return 255
	}

	return uint8(val)
}
