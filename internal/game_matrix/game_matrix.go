package gamematrix

import (
	"fmt"
	"slices"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/themilchenko/game_theory/internal/game_matrix/analytical"
	brownrobinson "github.com/themilchenko/game_theory/internal/game_matrix/brown_robinson"
	"gonum.org/v1/gonum/mat"
)

type GameMatrix struct {
	plainM [][]float64
	m      *mat.Dense

	lowestPrice  float64
	lowestIdx    int
	highestPrice float64
	highestIdx   int
}

func New(m [][]float64) (*GameMatrix, error) {
	if len(m) == 0 {
		return nil, fmt.Errorf("matrix should be not empty")
	}

	prev := len(m[0])

	for i := 1; i < len(m); i++ {
		if prev != len(m[i]) {
			return nil, fmt.Errorf("strings of matrix should be with the same size")
		}

		prev = len(m[i])
	}

	rows := len(m)
	cols := len(m[0])
	flatData := make([]float64, rows*cols)

	for i := range rows {
		for j := range cols {
			flatData[i*cols+j] = float64(m[i][j])
		}
	}

	g := &GameMatrix{
		plainM: m,
		m:      mat.NewDense(rows, cols, flatData),
	}
	g.lowestPrice, g.lowestIdx = g.calulateLowestPrice()
	g.highestPrice, g.highestIdx = g.calculateHighestPrice()

	return g, nil
}

func (g *GameMatrix) String() string {
	t := table.NewWriter()

	b := &strings.Builder{}
	t.SetOutputMirror(b)

	r, c := g.m.Dims()

	header := table.Row{"Strategies"}
	for i := range c {
		header = append(header, fmt.Sprintf("b_%d", i+1))
	}
	header = append(header, "min win of A player")
	t.AppendHeader(header)

	for i := range r {
		r := table.Row{fmt.Sprintf("a_%d", i+1)}

		for j := range c {
			r = append(r, fmt.Sprintf("%.3f", g.m.At(i, j)))
		}

		t.AppendRow(append(r, findMinInVec(g.m.RowView(i))))
	}

	t.AppendSeparator()

	raw := table.Row{"max loss of B player"}
	for j := range c {
		raw = append(raw, findMaxInVec(g.m.ColView(j)))
	}
	t.AppendRow(raw)

	s := table.StyleLight
	s.Format.Header = text.FormatLower
	t.SetStyle(s)

	t.Render()

	return b.String()
}

func (g *GameMatrix) MatrixString() string {
	return fmt.Sprintf("%.3v\n", mat.Formatted(g.m))
}

func (g *GameMatrix) LowestPrice() (float64, int) {
	return g.lowestPrice, g.lowestIdx
}

func (g *GameMatrix) HighestPrice() (float64, int) {
	return g.highestPrice, g.highestIdx
}

func (g *GameMatrix) SolveAnalytical() (*analytical.Solution, error) {
	solver, err := analytical.New(g.m)
	if err != nil {
		return nil, err
	}

	return solver.Solve()
}

func (g *GameMatrix) SolveBrownRobinson(opts ...brownrobinson.Opt) *brownrobinson.Solution {
	solver := brownrobinson.New(g.plainM, opts...)

	return solver.Solve()
}

func (g *GameMatrix) calulateLowestPrice() (float64, int) {
	minStrings := make([]float64, 0, g.m.RawMatrix().Rows)

	for i := range g.m.RawMatrix().Rows {
		minStrings = append(minStrings, findMinInVec(g.m.RowView(i)))
	}

	return slices.Max(minStrings), slices.Index(minStrings, slices.Max(minStrings))
}

func (g *GameMatrix) calculateHighestPrice() (float64, int) {
	maxColumns := make([]float64, 0, g.m.RawMatrix().Cols)

	for j := range g.m.RawMatrix().Cols {
		maxColumns = append(maxColumns, findMaxInVec(g.m.ColView(j)))
	}

	return slices.Min(maxColumns), slices.Index(maxColumns, slices.Min(maxColumns))
}
