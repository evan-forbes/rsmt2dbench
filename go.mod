module github.com/evan-forbes/rsmt2dbench

go 1.15

replace (
	github.com/lazyledger/lazyledger-core => /home/evan/go/src/github.com/lazyledger/lazyledger-core
	github.com/lazyledger/rsmt2d => /home/evan/go/src/github.com/lazyledger/rsmt2d
)

require (
	github.com/lazyledger/lazyledger-core v0.0.0-00010101000000-000000000000
	github.com/lazyledger/nmt v0.0.0-20201112204856-4bc77a77815c
	github.com/lazyledger/rsmt2d v0.0.0-20200626141417-ea94438fa457
	golang.org/x/crypto v0.0.0-20201117144127-c1f2f97bffc9
	gonum.org/v1/gonum v0.8.1
)