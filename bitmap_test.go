package index_test

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/espang/index"
)

var ages []int

func init() {
	f, err := os.Open("age.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	_ = scanner.Text()
	for scanner.Scan() {
		line := scanner.Text()
		v, err := strconv.ParseFloat(line, 64)
		if err != nil {
			continue
		}
		ages = append(ages, int(v))
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Println("Age has ", len(ages), " elements")
}

func TestBitmapIndex(t *testing.T) {
	idx := index.NewBitmapIndex(ages)

	fmt.Println("mem: ", idx.Size())
	fmt.Println("len: ", idx.Len())
	fmt.Println("uncompressed: ", idx.Len()*len(ages)/8)

	var total uint64
	for _, a := range ages {
		if a > 40 {
			total++
		}
	}
	start := time.Now()
	count, iter := idx.Greater(40)
	fmt.Printf("took %v\n", time.Since(start))
	fmt.Println(count, total)
	_ = iter

}

var gl int

func BenchmarkCreateBitmapIndex(b *testing.B) {
	var items int
	for n := 1; n < b.N; n++ {
		idx := index.NewBitmapIndex(ages)
		items = idx.Len()
	}
	gl = items
}
