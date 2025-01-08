package quantizer

import (
	"image"
	"image/color"
	"log"
	"sort"
)

type colorBucket struct {
	buff []colorVec
	axis int
}

func newBucket() colorBucket {
	return colorBucket{[]colorVec{}, 0}
}

func (cb *colorBucket) clear() {
	clear(cb.buff)
	cb.buff = cb.buff[:0]
}

func (cb *colorBucket) makeBucketBuff(src image.Image) {
	bounds := src.Bounds()
	for y := 0; y < bounds.Dy(); y++ {
		var colRow []colorVec
		for x := 0; x < bounds.Dx(); x++ {
			col := colorVecAt(src, x, y)
			cb.buff = append(cb.buff, col)
			colRow = append(colRow, col)
		}
	}
	cb.sort()
}

func (cb colorBucket) Len() int {
	return len(cb.buff)
}

func (cb colorBucket) Less(i, j int) bool {
	col1 := cb.buff[i]
	col2 := cb.buff[j]
	return col1[cb.axis] < col2[cb.axis]
}

func (cb colorBucket) Swap(i, j int) {
	cb.buff[i], cb.buff[j] = cb.buff[j], cb.buff[i]
}

func (cb *colorBucket) updateAxis() {
	cb.axis = span(cb.buff)
}

func (cb colorBucket) sort() {
	sort.Sort(cb)
}

func (cb colorBucket) part(start, end int) colorBucket {
	return colorBucket{cb.buff[start:end], cb.axis}
}

func (cb colorBucket) avgPart(start, end uint) color.RGBA {
	var avg [3]float32

	l := end - start
	sub := cb.buff[start:end]

	for axis := 0; axis < 3; axis++ {
		for _, col := range sub {
			avg[axis] += float32(col[axis]) / float32(l)
		}
	}

	return color.RGBA{uint8(avg[0]), uint8(avg[1]), uint8(avg[2]), 255}
}

func (cb colorBucket) getMrkers(depth int) (uint64, uint64) {
	nOfMarkers := uint64(1) << uint64(depth)
	if nOfMarkers == 0 {
		log.Fatal("nOfBuckets 0???")
	}

	partSize := uint64(cb.Len()) / nOfMarkers

	return nOfMarkers, partSize
}

func (cb colorBucket) getMarker(index, depth uint) uint {
	//return index * uint64(math.Round(float64(cb.Len())/float64(uint64(1)<<uint64(depth))))
	return index * uint(cb.Len()) / (uint(1) << depth)
}
