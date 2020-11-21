package erasure

import (
	"fmt"
	"testing"

	"github.com/lazyledger/rsmt2d"
)

// package level data dump to avoid unrealistic compiler optimization
var dump *rsmt2d.ExtendedDataSquare

/////////////////////////////////////////////
// Erasure Benchmarks
/////////////////////////////////////////////
func BenchmarkOriginalErasure(b *testing.B) {
	var res *rsmt2d.ExtendedDataSquare
	for i := uint(2); i < 50; i++ {
		data := genRandDS(i * i)
		b.Run(
			fmt.Sprintf("single thread data square width %d", i*i),
			func(b *testing.B) {
				for j := 0; j < b.N; j++ {
					eds := BaseScenario(data)
					res = eds
				}
			},
		)
	}
	dump = res
}

func BenchmarkParallelErasure(b *testing.B) {
	var res *rsmt2d.ExtendedDataSquare
	workers := 8
	for i := uint(2); i < 50; i++ {
		data := genRandDS(i * i)
		b.Run(
			fmt.Sprintf("multi-threaded %d thread data square width %d", workers, i*i),
			func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					eds := ParaScenario(data, 1)
					res = eds
				}
			},
		)
	}
	dump = res
}
