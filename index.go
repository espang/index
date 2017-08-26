package index

import (
	"bytes"
	"fmt"
	"sort"
)

type Index uint16

type Indexer interface {
	Get(val []byte, op Operator) Iter
}

type Iter interface {
	Next() bool
	Index() Index
	Err() error
}

type Operator int

const (
	Equal Operator = iota
	Less
	LessEqual
	Greater
	GreaterEqual
	NotEqual
)

type Uint16Iter interface {
	Next() bool
	Uint16() uint16
	Err() error
}

type index struct {
	values  [][]byte
	indices [][]Index
}

func (idx *index) add(val []byte, rowIdx Index) {
	if len(idx.values) == 0 {
		idx.values = [][]byte{val}
		idx.indices = [][]Index{[]Index{rowIdx}}
		return
	}

	for i, v := range idx.values {
		switch bytes.Compare(v, val) {
		case -1:
			// the value at index i is smaller than val
			continue
		case 0:
			idx.indices[i] = append(idx.indices[i], rowIdx)
			return
		case 1:
			idx.values = append(idx.values, nil)
			idx.indices = append(idx.indices, nil)

			copy(idx.values[i+1:], idx.values[i:])
			copy(idx.indices[i+1:], idx.indices[i:])

			idx.values[i] = val
			idx.indices[i] = []Index{rowIdx}
			return
		}
	}

	idx.values = append(idx.values, val)
	idx.indices = append(idx.indices, []Index{rowIdx})
}

func (idx *index) get(val []byte, op Operator) Iter {
	i := sort.Search(
		len(idx.values),
		func(i int) bool {
			return bytes.Compare(val, idx.values[i]) <= 0
		},
	)
	switch op {
	case Equal:
		if i < len(idx.values) && bytes.Equal(val, idx.values[i]) {
			return &equal{iterator{i: i, idx: idx}}
		}
		return null{}
	case NotEqual:
		return err{error: fmt.Errorf("not equal is not implemented")}
	case Less:
		return &left{iterator{i: i - 1, idx: idx}}
	case LessEqual:
		if i < len(idx.values) && bytes.Equal(val, idx.values[i]) {
			return &left{iterator{i: i, idx: idx}}
		}
		return &left{iterator{i: i - 1, idx: idx}}
	case Greater:
		return err{error: fmt.Errorf("greater is not implemented")}
	case GreaterEqual:
		return err{error: fmt.Errorf("not greater is not implemented")}
	default:
		return err{error: fmt.Errorf("unknown operator: %d", op)}
	}
}

type iterator struct {
	i, j int
	idx  *index
}

func (i *iterator) Index() Index { return i.idx.indices[i.i][i.j-1] }
func (i *iterator) Err() error   { return nil }

type null struct{}

func (null) Next() bool   { return false }
func (null) Index() Index { return 0 }
func (null) Err() error   { return nil }

type err struct {
	null
	error
}

func (e err) Err() error { return e.error }

type equal struct {
	iterator
}

func (e *equal) Next() bool {
	e.j++
	return e.j <= len(e.idx.indices[e.i])
}

type right struct {
	iterator
}

func (i *right) Next() bool {
	if i.j == len(i.idx.indices[i.i]) {
		i.j = 1
		i.i++
		return i.i < len(i.idx.indices)
	}
	i.j++
	return i.i < len(i.idx.indices)
}

type left struct {
	iterator
}

func (i *left) Next() bool {
	if i.i == -1 {
		return false
	}
	if i.j == len(i.idx.indices[i.i]) {
		i.j = 1
		i.i--
		return i.i >= 0
	}
	i.j++
	return i.i >= 0
}
