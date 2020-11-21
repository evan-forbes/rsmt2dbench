package daheader

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	types "github.com/lazyledger/lazyledger-core/types"
	"github.com/lazyledger/nmt"
	"github.com/lazyledger/nmt/namespace"
	"github.com/lazyledger/rsmt2d"
	"golang.org/x/crypto/sha3"
	"gonum.org/v1/gonum/stat"
)

const NamespaceSize = 8

var (
	hash                    []byte // package level var to disable unrealistic compiler optimization
	newBaseHashFunc         = sha3.New256
	ParitySharesNamespaceID = namespace.ID{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
)

func RunCombinedBenchmarks() {
	// benchmark the single threaded implementation
	for i := 4; i < 129; i += 4 {
		fmt.Println(BenchSingleDAHash(20, i))
	}
	fmt.Println("\n------------------------------------\n")
	for w := 1; w < 257; w *= 2 {
		for i := 4; i < 129; i += 4 {
			fmt.Println(BenchMultiDAHash(20, i, w))
		}
		fmt.Println("\n------------------------------------\n")
	}
}

func BenchSingleDAHash(count, width int) result {
	var times []uint64
	shares, raw := MockShares(width)
	for i := 0; i < count; i++ {
		start := time.Now().UnixNano()
		eds, err := rsmt2d.ComputeExtendedDataSquare(raw, rsmt2d.RSGF8)
		if err != nil {
			panic(fmt.Sprintf("unexpected error: %v", err))
		}
		HashOriginal(eds, shares, 0)
		end := time.Now().UnixNano()
		elapsed := end - start
		times = append(times, uint64(elapsed))
	}
	return averageResult(0, width, false, times)
}

// BenchMultiDaHash runs a benchmark to determine the
func BenchMultiDAHash(count, width int, workers int) result {
	var times []uint64
	shares, raw := MockShares(width)
	for i := 0; i < count; i++ {
		start := time.Now().UnixNano()
		eds, err := rsmt2d.ParallelComputeExtendedDataSquare(raw, rsmt2d.NewRSGF8Codec(), workers)
		if err != nil {
			panic(fmt.Sprintf("unexpected error: %v", err))
		}
		HashMultiThread(eds, shares, workers)
		end := time.Now().UnixNano()
		elapsed := end - start
		times = append(times, uint64(elapsed))
	}
	return averageResult(workers, width, true, times)
}

func RunNMTBenchmarks() {
	// benchmark the single threaded implementation
	for i := 4; i < 129; i += 4 {
		fmt.Println(BenchNMTGeneration(40, i, HashOriginal, 0))
	}
	fmt.Println("\n------------------------------------\n")
	for w := 1; w < 257; w *= 2 {
		for i := 4; i < 129; i += 4 {
			fmt.Println(BenchNMTGeneration(40, i, HashMultiThread, w))
		}
		fmt.Println("\n------------------------------------\n")
	}
}

type hashBench func(eds *rsmt2d.ExtendedDataSquare, nss types.NamespacedShares, workers int)

func BenchNMTGeneration(count, width int, f hashBench, workers int) result {
	var times []uint64
	shares, raw := MockShares(width)
	eds, err := rsmt2d.ComputeExtendedDataSquare(raw, rsmt2d.RSGF8)
	if err != nil {
		panic(fmt.Sprintf("unexpected error: %v", err))
	}
	for i := 0; i < count; i++ {
		start := time.Now().UnixNano()
		f(eds, shares, workers)
		end := time.Now().UnixNano()
		elapsed := end - start
		times = append(times, uint64(elapsed))
	}
	return averageResult(workers, width, workers > 1, times)
}

// HashOriginal runs the current implementation to generate nmts. This function
// moves the unexported version from lazyledger-core the unexported version and
// copies it here
func HashOriginal(eds *rsmt2d.ExtendedDataSquare, nss types.NamespacedShares, workers int) {
	// compute roots:
	squareWidth := eds.Width()
	originalDataWidth := squareWidth / 2
	dah := types.DataAvailabilityHeader{
		RowsRoots:   make([]namespace.IntervalDigest, squareWidth),
		ColumnRoots: make([]namespace.IntervalDigest, squareWidth),
	}

	// compute row and column roots:
	for outerIdx := uint(0); outerIdx < squareWidth; outerIdx++ {
		rowTree := nmt.New(newBaseHashFunc(), nmt.NamespaceIDSize(NamespaceSize))
		colTree := nmt.New(newBaseHashFunc(), nmt.NamespaceIDSize(NamespaceSize))
		for innerIdx := uint(0); innerIdx < squareWidth; innerIdx++ {
			if outerIdx < originalDataWidth && innerIdx < originalDataWidth {
				mustPush(rowTree, nss[outerIdx*originalDataWidth+innerIdx])
				mustPush(colTree, nss[innerIdx*originalDataWidth+outerIdx])
			} else {
				rowData := eds.Row(outerIdx)
				colData := eds.Column(outerIdx)

				parityCellFromRow := rowData[innerIdx]
				parityCellFromCol := colData[innerIdx]
				// FIXME(ismail): do not hardcode usage of PrefixedData8 here:
				mustPush(rowTree, namespace.PrefixedData8(
					append(ParitySharesNamespaceID, parityCellFromRow...),
				))
				mustPush(colTree, namespace.PrefixedData8(
					append(ParitySharesNamespaceID, parityCellFromCol...),
				))
			}
		}
		dah.RowsRoots[outerIdx] = rowTree.Root()
		dah.ColumnRoots[outerIdx] = colTree.Root()
	}

	hash = dah.Hash()
}

// HashMultiThread runs the multithreaded implementation to generate nmts for the
// dah. This function moves the unexported version from lazyledger-core the
// unexported version and copies it here
func HashMultiThread(eds *rsmt2d.ExtendedDataSquare, nss types.NamespacedShares, workers int) {
	squareWidth := eds.Width()
	originalDataWidth := squareWidth / 2
	dah := types.DataAvailabilityHeader{
		Mut:         sync.Mutex{},
		RowsRoots:   make([]namespace.IntervalDigest, squareWidth),
		ColumnRoots: make([]namespace.IntervalDigest, squareWidth),
	}
	work := func(wg *sync.WaitGroup, dah *types.DataAvailabilityHeader, jobs <-chan uint) {
		defer wg.Done()
		for outerIdx := range jobs {
			rowTree := nmt.New(newBaseHashFunc(), nmt.NamespaceIDSize(NamespaceSize))
			colTree := nmt.New(newBaseHashFunc(), nmt.NamespaceIDSize(NamespaceSize))
			for innerIdx := uint(0); innerIdx < squareWidth; innerIdx++ {
				if outerIdx < originalDataWidth && innerIdx < originalDataWidth {
					mustPush(rowTree, nss[outerIdx*originalDataWidth+innerIdx])
					mustPush(colTree, nss[innerIdx*originalDataWidth+outerIdx])
				} else {
					rowData := eds.Row(outerIdx)
					colData := eds.Column(outerIdx)

					parityCellFromRow := rowData[innerIdx]
					parityCellFromCol := colData[innerIdx]
					// FIXME(ismail): do not hardcode usage of PrefixedData8 here:
					mustPush(rowTree, namespace.PrefixedData8(
						append(ParitySharesNamespaceID, parityCellFromRow...),
					))
					mustPush(colTree, namespace.PrefixedData8(
						append(ParitySharesNamespaceID, parityCellFromCol...),
					))
				}
			}

			rowRoot := rowTree.Root()
			colRoot := colTree.Root()
			dah.Mut.Lock()
			dah.RowsRoots[outerIdx] = rowRoot
			dah.ColumnRoots[outerIdx] = colRoot
			dah.Mut.Unlock()
		}
	}

	jobs := make(chan uint, squareWidth)
	go func() {
		defer close(jobs)
		for id := uint(0); id < squareWidth; id++ {
			jobs <- id
		}
	}()
	wg := &sync.WaitGroup{}
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go work(wg, &dah, jobs)
	}

	wg.Wait()

	hash = dah.Hash()
}

func mustPush(rowTree *nmt.NamespacedMerkleTree, namespacedShare namespace.Data) {
	if err := rowTree.Push(namespacedShare); err != nil {
		panic(
			fmt.Sprintf("invalid data; could not push share to tree: %#v, err: %v",
				namespacedShare,
				err,
			),
		)
	}
}

// MockShares makes random data in 256byte shares of length `width`
func MockShares(width int) (shares types.NamespacedShares, rawShares [][]byte) {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	count := width * width
	for i := 0; i < count; i++ {
		// generate random data
		s := make([]byte, 256)
		r.Read(s)
		share := types.NamespacedShare{
			Share: s,
			ID:    mockID(i),
		}
		shares = append(shares, share)
		rawShares = append(rawShares, s)
	}
	return shares, rawShares
}

func mockID(id int) namespace.ID {
	switch {
	case id < 256:
		return namespace.ID{0, 0, 0, 0, 0, 0, 0, byte(id)}
	default:
		return namespace.ID{0, 0, 0, 0, 0, 0, 100, 0}
	}
}

type result struct {
	workers  int
	threaded bool
	time     uint64
	dev      uint
	square   int
}

func (r result) String() string {
	return fmt.Sprintf("%d %d %d", r.square, r.workers, r.time)
}

func averageResult(workers, square int, threaded bool, input []uint64) result {
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
		square:   square,
	}
}
