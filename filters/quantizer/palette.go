package quantizer

import (
	"image/color"
)

type Palette interface {
	conv(c color.RGBA) color.RGBA
	convVec(c colorVec) color.RGBA
	palette() color.Palette
}

type MemoPalette struct {
	StdPalette
	Memo map[color.RGBA]color.RGBA
}

func newMemo(p StdPalette) MemoPalette {
	return MemoPalette{p, make(map[color.RGBA]color.RGBA)}
}

func (p MemoPalette) conv(c color.RGBA) color.RGBA {
	if col, ok := p.Memo[c]; !ok {
		qCol := p.StdPalette.conv(c)
		p.Memo[c] = qCol

		return qCol
	} else {

		return col
	}
}

func (p MemoPalette) convVec(cv colorVec) color.RGBA {
	c := color.RGBA{cv[0], cv[1], cv[2], 255}

	if col, ok := p.Memo[c]; !ok {
		qCol := p.StdPalette.conv(c)
		p.Memo[c] = qCol

		return qCol
	} else {

		return col
	}
}

type StdPalette color.Palette

func (p StdPalette) palette() color.Palette {
	return color.Palette(p)
}

func (p StdPalette) conv(col color.RGBA) color.RGBA {
	return color.Palette(p).Convert(col).(color.RGBA)
}

func (p StdPalette) convVec(col colorVec) color.RGBA {
	return color.Palette(p).Convert(
		color.RGBA{col[0], col[1], col[2], 255},
	).(color.RGBA)
}

type MapPalette map[colorVec]color.RGBA

func (p MapPalette) palette() color.Palette {
	var pal color.Palette
	for _, c := range p {
		pal = append(pal, c)
	}
	return pal
}

func (p MapPalette) conv(col color.RGBA) color.RGBA {
	return p.safeGet(colorVec{col.R, col.G, col.B})
}

func (p MapPalette) convVec(col colorVec) color.RGBA {
	return p.safeGet(col)
}

func (p MapPalette) safeGet(col colorVec) color.RGBA {
	if col, ok := p[col]; !ok {
		return color.RGBA{0, 0, 0, 255}
	} else {
		return col
	}
}
