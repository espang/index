package index

import (
	"crypto/rand"
	"reflect"
	"testing"
)

var (
	arrayLength64kWidth1 [][]byte
	arrayLength64kWidth2 [][]byte
)

func init() {
	arrayLength64kWidth1 = randArray(1, 1<<16)
	arrayLength64kWidth2 = randArray(2, 1<<16)
}

var g int

func BenchmarkCreateIndex(b *testing.B) {
	var items int
	for n := 1; n < b.N; n++ {
		idx := &index{}
		for i, v := range arrayLength64kWidth1 {
			idx.add(v, Index(i))
		}
		items = len(idx.values)
	}
	g = items
}

func randArray(width, length int) [][]byte {
	arr := make([][]byte, length)
	for i := range arr {
		b := make([]byte, width)
		_, _ = rand.Read(b)
		arr[i] = b
	}
	return arr
}

func consume(iter Iter) []Index {
	var res []Index
	for iter.Next() {
		res = append(res, iter.Index())
	}
	return res
}

func TestIndex(t *testing.T) {

	idx := &index{}

	idx.add([]byte{1}, 0)
	idx.add([]byte{2}, 1)
	idx.add([]byte{1}, 2)
	idx.add([]byte{4}, 3)
	idx.add([]byte{1}, 4)
	idx.add([]byte{2}, 5)
	idx.add([]byte{0}, 6)

	values := [][]byte{
		{0}, {1}, {2}, {4},
	}
	if got := idx.values; !reflect.DeepEqual(got, values) {
		t.Errorf("got %v; want %v", got, values)
	}

	indices := [][]Index{
		{6}, {0, 2, 4}, {1, 5}, {3},
	}
	if got := idx.indices; !reflect.DeepEqual(got, indices) {
		t.Errorf("got %v; want %v", got, indices)
	}

	testCases := []struct {
		name string
		val  []byte
		want []Index
		op   Operator
	}{
		{"< 0", []byte{0}, nil, Less},
		{"< 1", []byte{1}, []Index{6}, Less},
		{"< 2", []byte{2}, []Index{0, 2, 4, 6}, Less},
		{"< 3", []byte{3}, []Index{1, 5, 0, 2, 4, 6}, Less},
		{"< 4", []byte{4}, []Index{1, 5, 0, 2, 4, 6}, Less},
		{"< 5", []byte{5}, []Index{3, 1, 5, 0, 2, 4, 6}, Less},

		{"<= 0", []byte{0}, []Index{6}, LessEqual},
		{"<= 1", []byte{1}, []Index{0, 2, 4, 6}, LessEqual},
		{"<= 2", []byte{2}, []Index{1, 5, 0, 2, 4, 6}, LessEqual},
		{"<= 3", []byte{3}, []Index{1, 5, 0, 2, 4, 6}, LessEqual},
		{"<= 4", []byte{4}, []Index{3, 1, 5, 0, 2, 4, 6}, LessEqual},
		{"<= 5", []byte{5}, []Index{3, 1, 5, 0, 2, 4, 6}, LessEqual},

		{"= 0", []byte{0}, []Index{6}, Equal},
		{"= 1", []byte{1}, []Index{0, 2, 4}, Equal},
		{"= 2", []byte{2}, []Index{1, 5}, Equal},
		{"= 3", []byte{3}, nil, Equal},
		{"= 4", []byte{4}, []Index{3}, Equal},
		{"= 5", []byte{5}, nil, Equal},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := consume(idx.get(tc.val, tc.op))
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %v; want %v", got, tc.want)
			}
		})
	}
}
