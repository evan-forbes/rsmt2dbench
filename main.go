package main

import (
	"fmt"

	"github.com/evan-forbes/rsmt2dbench/daheader"
	"github.com/evan-forbes/rsmt2dbench/erasure"
	"gonum.org/v1/gonum/stat"
)

func main() {
	// shares, raw := daheader.MockShares(128)
	// eds, err := rsmt2d.ComputeExtendedDataSquare(raw, rsmt2d.RSGF8)
	// if err != nil {
	// 	panic(fmt.Sprintf("unexpected error: %v", err))
	// }
	// trace.Start(os.Stdout)
	// defer trace.Stop()
	// daheader.HashMultiThread(eds, shares, 16)

	daheader.RunCombinedBenchmarks()

}

type result struct {
	workers  int
	threaded bool
	time     uint64
	dev      uint
}

func (r result) String() string {
	return fmt.Sprintf("%d %d (+/-) %d", r.workers, r.time, r.dev)
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

func runAveragedRSMT2DBenchmarks() {
	// for i := uint(4); i < 129; i += 4 {
	// 	var res uint64
	// 	var results []uint64
	// 	for j := 0; j < 30; j++ {
	// 		res = erasure.BenchSingleThread(i)
	// 		results = append(results, res)
	// 	}
	// 	final := averageResult(1, false, results)
	// 	fmt.Printf("%d %d %t %d\n", i, final.workers, final.threaded, final.time)
	// }
	// fmt.Println("\n---------------------------------\n")
	// run the multithreaded benchmarks
	// for i := uint(2); i < 128; i++ {
	// 	var res uint64
	// 	var results []uint64
	// 	for j := 0; j < 3; j++ {
	// 		res = erasure.BenchMultiThreaded(i, workers)
	// 		results = append(results, res)
	// 	}
	// 	final := averageResult(workers, true, results)
	// 	fmt.Printf("%d %d %t %d +/- %d\n", i, final.workers, final.threaded, final.time, final.dev)
	// }
	for k := 64; k < 129; k = k * 2 {
		for i := uint(4); i < 129; i += 4 {
			results := []uint64{}
			// run the multithreaded implementation
			for j := 0; j < 40; j++ {
				res := erasure.BenchMultiThreaded(i, k)
				results = append(results, res)
			}
			final := averageResult(k, false, results)
			fmt.Printf("%d %d %d\n", i, final.workers, final.time)
		}
		fmt.Println("\n---------------------------------\n")
	}
}
