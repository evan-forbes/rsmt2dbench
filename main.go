package main

import (
	"crypto/rand"
	"os"
	"runtime/trace"
	"time"

	"github.com/lazyledger/rsmt2d"
	"gonum.org/v1/gonum/stat"
)

func main() {
	trace.Start(os.Stdout)
	defer trace.Stop()
	data := genRandDS(625)
	ParaScenario(data, 8)
	// script to average the benchmarks for 30 runs a piece.
	// the typical go benchmarks we're giving me identical results across the board

	// for i := uint(2); i < 51; i++ {
	// 	var res uint64
	// 	var results []uint64
	// 	for j := 0; j < 30; j++ {
	// 		res = BenchSingleThread(i)
	// 		results = append(results, res)
	// 	}
	// 	final := averageResult(1, false, results)
	// 	fmt.Printf("%d %d %t %d\n", i*i, final.workers, final.threaded, final.time)
	// }
	// run the multithreaded benchmarks
	// workers := 32
	// for i := uint(2); i < 51; i++ {
	// 	var res uint64
	// 	var results []uint64
	// 	for j := 0; j < 50; j++ {
	// 		res = BenchMultiThreaded(i, workers)
	// 		results = append(results, res)
	// 	}
	// 	final := averageResult(workers, true, results)
	// 	fmt.Printf("%d %d %t %d +/- %d\n", i*i, final.workers, final.threaded, final.time, final.dev)
	// }
	// for k := 2; k < 33; k = k * 2 {
	// 	results = []uint64{}
	// 	for i := uint(2); i < 50; i++ {
	// 		// run the multithreaded implementation
	// 		res = BenchMultiThreaded(i, k)
	// 		results = append(results, res)
	// 	}
	// 	final = averageResult(1, false, results)
	// 	fmt.Printf("%d %t %d", final.workers, final.threaded, final.time)
	// }

}

type result struct {
	workers  int
	threaded bool
	time     uint64
	dev      uint
}

func averageResult(workers int, threaded bool, input []uint64) result {
	var total uint64
	var floatIn []float64
	for _, i := range input {
		total += i
		floatIn = append(floatIn, float64(i))
	}

	return result{
		workers:  workers,
		threaded: threaded,
		time:     total / uint64(len(input)),
		dev:      uint(stat.StdDev(floatIn, nil)),
	}
}

func BenchSingleThread(i uint) uint64 {
	data := genRandDS(i * i)
	start := time.Now().UnixNano()
	BaseScenario(data)
	end := time.Now().UnixNano()
	elapsed := end - start
	return uint64(elapsed)
}

func BenchMultiThreaded(i uint, workers int) uint64 {
	data := genRandDS(i * i)
	start := time.Now().UnixNano()
	ParaScenario(data, workers)
	end := time.Now().UnixNano()
	elapsed := end - start
	return uint64(elapsed)
}

// BaseScenario s
func BaseScenario(data [][]byte) (eds *rsmt2d.ExtendedDataSquare) {
	eds, err := rsmt2d.ComputeExtendedDataSquare(data, rsmt2d.RSGF8)
	if err != nil {
		panic(err)
	}
	return eds
}

// use the default source of psuedo randomness to generation a data square of
// size width^2
func genRandDS(width uint) [][]byte {
	var ds [][]byte
	for i := uint(0); i < width; i++ {
		row := make([]byte, width)
		rand.Read(row)
		ds = append(ds, row)
	}
	return ds
}

func ParaScenario(data [][]byte, workers int) (eds *rsmt2d.ExtendedDataSquare) {
	eds, err := rsmt2d.ParallelComputeExtendedDataSquare(data, rsmt2d.NewRSGF8Codec(), workers)
	if err != nil {
		panic(err)
	}
	return eds
}
