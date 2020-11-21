package main

import (
	core "github.com/lazyledger/lazyledger-core/types"
	"github.com/lazyledger/nmt"
	"github.com/lazyledger/nmt/namespace"
	"github.com/lazyledger/rsmt2d"
)

// package level var to disable unrealistic compiler optimization
var hash string

newBaseHashFunc = sha3.New256

func HashOriginal(eds *rsmt2d.ExtendedDataSquare) {
	// compute roots:
	squareWidth := eds.Width()
	originalDataWidth := squareWidth / 2
	b.DataAvailabilityHeader = core.DataAvailabilityHeader{
		RowsRoots:   make([]namespace.IntervalDigest, squareWidth),
		ColumnRoots: make([]namespace.IntervalDigest, squareWidth),
	}

	// compute row and column roots:
	for outerIdx := uint(0); outerIdx < squareWidth; outerIdx++ {
		rowTree := nmt.New(newBaseHashFunc(), nmt.NamespaceIDSize(NamespaceSize))
		colTree := nmt.New(newBaseHashFunc(), nmt.NamespaceIDSize(NamespaceSize))
		for innerIdx := uint(0); innerIdx < squareWidth; innerIdx++ {
			if outerIdx < originalDataWidth && innerIdx < originalDataWidth {
				mustPush(rowTree, namespacedShares[outerIdx*originalDataWidth+innerIdx])
				mustPush(colTree, namespacedShares[innerIdx*originalDataWidth+outerIdx])
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
		b.DataAvailabilityHeader.RowsRoots[outerIdx] = rowTree.Root()
		b.DataAvailabilityHeader.ColumnRoots[outerIdx] = colTree.Root()
	}

	b.DataHash = b.DataAvailabilityHeader.Hash()
}
