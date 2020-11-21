package erasure

import (
	"crypto/rand"
	"time"

	"github.com/lazyledger/rsmt2d"
)

// BenchSingleThread runs the homecooked benchmark for the original implemenation
func BenchSingleThread(i uint) uint64 {
	data := genRandDS(i)
	start := time.Now().UnixNano()
	BaseScenario(data)
	end := time.Now().UnixNano()
	elapsed := end - start
	return uint64(elapsed)
}

// BenchMultiThreaded runs the homecooked multithreaded benchmark
func BenchMultiThreaded(i uint, workers int) uint64 {
	data := genRandDS(i)
	start := time.Now().UnixNano()
	ParaScenario(data, workers)
	end := time.Now().UnixNano()
	elapsed := end - start
	return uint64(elapsed)
}

// BaseScenario runs the original implementation
func BaseScenario(data [][]byte) (eds *rsmt2d.ExtendedDataSquare) {
	eds, err := rsmt2d.ComputeExtendedDataSquare(data, rsmt2d.RSGF8)
	if err != nil {
		panic(err)
	}
	return eds
}

// ParaScenario runs the multithreaded implemenation
func ParaScenario(data [][]byte, workers int) (eds *rsmt2d.ExtendedDataSquare) {
	eds, err := rsmt2d.ParallelComputeExtendedDataSquare(data, rsmt2d.NewRSGF8Codec(), workers)
	if err != nil {
		panic(err)
	}
	return eds
}

// // use the default source of psuedo randomness to generation a data square of
// // size width^2
// func genRandDS(width uint) [][]byte {
// 	var ds [][]byte
// 	for i := uint(0); i < width; i++ {
// 		row := make([]byte, width)
// 		rand.Read(row)
// 		ds = append(ds, row)
// 	}
// 	return ds
// }

// genRandDS make a datasquare of random data, with width describing the number
// of shares on a single side of the ds
func genRandDS(width uint) [][]byte {
	var ds [][]byte
	count := width * width
	for i := uint(0); i < count; i++ {
		share := make([]byte, 256)
		rand.Read(share)
		ds = append(ds, share)
	}
	return ds
}
