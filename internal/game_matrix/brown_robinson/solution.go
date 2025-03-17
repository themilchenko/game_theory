package brownrobinson

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

type Solution struct {
	t table.Writer
	b *strings.Builder

	x []float64
	y []float64
	v float64
}

func newSolution(xLen, yLen int) *Solution {
	s := &Solution{
		t: table.NewWriter(),
		b: &strings.Builder{},
		x: make([]float64, 0),
		y: make([]float64, 0),
	}
	s.t.SetOutputMirror(s.b)

	header := table.Row{"#", "A", "B"}

	for i := range xLen {
		header = append(header, fmt.Sprintf("x_%d", i+1))
	}
	for i := range yLen {
		header = append(header, fmt.Sprintf("y_%d", i+1))
	}

	header = append(header, "top_game_cost", "lower_game_cost", "eps")

	s.t.AppendHeader(header)

	st := table.StyleLight
	st.Format.Header = text.FormatLower
	s.t.SetStyle(st)

	return s
}

func (s *Solution) String() string {
	xStr := &strings.Builder{}
	xStr.WriteString("x* = (")
	for i, v := range s.x {
		fmt.Fprintf(xStr, "%.3f", v)

		if i != len(s.x)-1 {
			xStr.WriteString(", ")
		}
	}

	xStr.WriteString(")\n")

	s.b.WriteString(xStr.String())

	yStr := &strings.Builder{}
	yStr.WriteString("y* = (")
	for i, v := range s.y {
		fmt.Fprintf(yStr, "%.3f", v)

		if i != len(s.y)-1 {
			yStr.WriteString(", ")
		}
	}

	yStr.WriteString(")\n")

	s.b.WriteString(yStr.String())

	fmt.Fprintf(s.b, "v = %.3f", s.v)

	return s.b.String()
}

func (s *Solution) append(it iter) {
	r := table.Row{it.num, fmt.Sprintf("x_%d", it.x+1),
		fmt.Sprintf("y_%d", it.y+1)}

	for _, v := range it.aWin {
		r = append(r, v)
	}

	for _, v := range it.bLoss {
		r = append(r, v)
	}

	r = append(r, fmt.Sprintf("%.3f", it.top),
		fmt.Sprintf("%.3f", it.lower),
		fmt.Sprintf("%.3f", it.eps))

	s.t.AppendRow(r)
}

func (s *Solution) finish(iters []iter) {
	for i := range iters[0].aWin {
		cnt := float64(0)

		for _, v := range iters {
			if v.x == i {
				cnt++
			}
		}

		s.x = append(s.x, cnt/float64(len(iters)))
	}

	for i := range iters[0].bLoss {
		cnt := float64(0)

		for _, v := range iters {
			if v.y == i {
				cnt++
			}
		}

		s.y = append(s.y, cnt/float64(len(iters)))
	}

	min := iters[0].top
	for i := 1; i < len(iters); i++ {
		if min > iters[i].top {
			min = iters[i].top
		}
	}

	max := iters[0].lower
	for i := 1; i < len(iters); i++ {
		if max < iters[i].lower {
			max = iters[i].lower
		}
	}

	s.v = (min + max) / 2

	s.t.Render()

	s.drawGraphics(iters)
}

func (s *Solution) drawGraphics(iters []iter) {
	p := plot.New()

	p.Title.Text = "Graph of convergence of upper and lower " +
		"game prices in the Brown-Robinson algorithm"
	p.X.Label.Text = "Iterations"
	p.Y.Label.Text = "Costs of game"

	lower := make([]float64, 0, len(iters))
	for _, v := range iters {
		lower = append(lower, v.lower)
	}
	lowerXYs := s.getCostPoints(lower)

	top := make([]float64, 0, len(iters))
	for _, v := range iters {
		top = append(top, v.top)
	}
	topXYs := s.getCostPoints(top)

	plotutil.AddLinePoints(p, "Lower cost", lowerXYs, "Top cost", topXYs)

	p.Save(10*vg.Inch, 5*vg.Inch, filepath.Join("artifacts", "lw1", "costs.png"))

	p = plot.New()
	p.Title.Text = "Estimation graph of the Brown-Robinson algorithm"
	p.X.Label.Text = "Iterations"
	p.Y.Label.Text = "Epsilon"

	eps := make([]float64, 0, len(iters))
	for _, v := range iters {
		eps = append(eps, v.eps)
	}
	epsXYs := s.getCostPoints(eps)

	plotutil.AddLinePoints(p, "Eps", epsXYs)

	p.Save(10*vg.Inch, 5*vg.Inch, filepath.Join("artifacts", "lw1", "estimation.png"))
}

func (s *Solution) getCostPoints(costs []float64) plotter.XYs {
	pts := make(plotter.XYs, len(costs))

	for i, v := range costs {
		pts[i].X = float64(i)
		pts[i].Y = v
	}

	return pts
}
