package brownrobinson

import (
	"math/rand"
	"slices"
)

const epsilon float64 = 0.1

type BrownRobinson struct {
	m     [][]float64
	iters []iter

	minTop   float64
	maxLower float64

	sol *Solution
}

type iter struct {
	num int

	aWin  []float64
	bLoss []float64

	x int
	y int

	top   float64
	lower float64
	eps   float64
}

func New(m [][]float64) *BrownRobinson {
	br := &BrownRobinson{
		m: m,
	}

	init := iter{
		num:   1,
		aWin:  br.column(0),
		bLoss: br.row(0),
		x:     0,
		y:     0,
	}

	init.top = slices.Max(init.aWin)
	init.lower = slices.Min(init.bLoss)
	init.eps = init.top - init.lower

	br.iters = []iter{init}

	br.sol = newSolution(len(init.aWin), len(init.bLoss))

	br.sol.append(init)

	br.minTop = init.top
	br.maxLower = init.lower

	return br
}

func (b *BrownRobinson) Solve() *Solution {
	// Start with 2 because first one in New() constructor.
	iterNum := 2

	for b.iters[len(b.iters)-1].eps > epsilon {
		b.step(iterNum)

		b.sol.append(b.iters[len(b.iters)-1])

		iterNum++
	}

	b.sol.finish(b.iters)

	return b.sol
}

func (b *BrownRobinson) step(iterNum int) {
	last := b.iters[len(b.iters)-1]

	it := iter{
		num: iterNum,
	}

	_, x := b.min(last.bLoss)
	_, y := b.max(last.aWin)

	it.x = x
	it.y = y

	it.aWin = make([]float64, 0, len(last.aWin))
	for i, v := range b.column(x) {
		it.aWin = append(it.aWin, last.aWin[i]+v)
	}

	it.bLoss = make([]float64, 0, len(last.bLoss))
	for i, v := range b.row(y) {
		it.bLoss = append(it.bLoss, last.bLoss[i]+v)
	}

	v, _ := b.max(it.aWin)
	it.top = v / float64(iterNum)

	if b.minTop > it.top {
		b.minTop = it.top
	}

	v, _ = b.min(it.bLoss)
	it.lower = v / float64(iterNum)

	if b.maxLower < it.lower {
		b.maxLower = it.lower
	}

	it.eps = b.minTop - b.maxLower

	b.iters = append(b.iters, it)
}

func (b *BrownRobinson) column(j int) []float64 {
	col := make([]float64, 0, len(b.m))

	for i := range b.m {
		col = append(col, b.m[i][j])
	}

	return col
}

func (b *BrownRobinson) row(i int) []float64 {
	row := make([]float64, 0, len(b.m[i]))

	return append(row, b.m[i]...)
}

func (b *BrownRobinson) max(s []float64) (float64, int) {
	mIdxs := make([]int, 0)

	m := slices.Max(s)

	for i, v := range s {
		if v == m {
			mIdxs = append(mIdxs, i)
		}
	}

	random := rand.Intn(len(mIdxs))

	return m, mIdxs[random]
}

func (b *BrownRobinson) min(s []float64) (float64, int) {
	mIdxs := make([]int, 0)

	m := slices.Min(s)

	for i, v := range s {
		if v == m {
			mIdxs = append(mIdxs, i)
		}
	}

	random := rand.Intn(len(mIdxs))

	return m, mIdxs[random]
}
