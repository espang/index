package index

import (
	"sort"

	"github.com/RoaringBitmap/roaring"
)

type BitmapIndex struct {
	//sorted list of containing values:
	values  []int
	bitmaps []*roaring.Bitmap
}

func (idx BitmapIndex) Len() int { return len(idx.values) }

func (idx BitmapIndex) Size() uint64 {
	var total uint64
	for _, bm := range idx.bitmaps {
		total += bm.GetSerializedSizeInBytes()
	}
	return total
}

func NewBitmapIndex(values []int) *BitmapIndex {
	m := map[int]*roaring.Bitmap{}
	for i, v := range values {
		if _, ok := m[v]; !ok {
			m[v] = roaring.NewBitmap()
		}
		m[v].AddInt(i)
	}

	vals := make([]int, 0, len(m))
	for i := range m {
		vals = append(vals, i)
	}

	sort.Ints(vals)
	bms := make([]*roaring.Bitmap, len(m))

	for i, v := range vals {
		bms[i] = m[v]
	}

	return &BitmapIndex{
		values:  vals,
		bitmaps: bms,
	}
}

func (idx *BitmapIndex) Greater(val int) (uint64, roaring.IntIterable) {
	i := sort.Search(
		len(idx.values),
		func(i int) bool { return idx.values[i] > val },
	)
	bm := roaring.ParOr(4, idx.bitmaps[i:]...)
	return bm.GetCardinality(), bm.Iterator()
}
