package main

import (
	"time"

	"github.com/lazyledger/rsmt2d"
)

// BenchSingleThread runs the homecooked benchmark for the original implemenation
func BenchSingleThread(i uint) uint64 {
	data := genRandDS(i * i)
	start := time.Now().UnixNano()
	BaseScenario(data)
	end := time.Now().UnixNano()
	elapsed := end - start
	return uint64(elapsed)
}

// BenchMultiThreaded runs the homecooked multithreaded benchmark
func BenchMultiThreaded(i uint, workers int) uint64 {
	data := genRandDS(i * i)
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
