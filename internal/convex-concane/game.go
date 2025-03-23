package convexconcane

import (
	"errors"
	"fmt"
	"math"
	"slices"
	"strings"

	gamematrix "github.com/themilchenko/game_theory/internal/game_matrix"
	brownrobinson "github.com/themilchenko/game_theory/internal/game_matrix/brown_robinson"
)

const (
	kernelFunc = "%.2fx^2 + %.2fy^2 + %.2fxy + %.2fx + %.2fy\n"
	eps        = 0.01
	lastIters  = 4
)

type ConvexConcane struct {
	a float64
	b float64
	c float64
	d float64
	e float64
}

func New(a, b, c, d, e float64) (*ConvexConcane, error) {
	if 2*a >= 0 || 2*b <= 0 {
		return nil, errors.New("game is not convex-concane")
	}

	return &ConvexConcane{
		a: a,
		b: b,
		c: c,
		d: d,
		e: e,
	}, nil
}

func (c *ConvexConcane) h(x, y float64) float64 {
	return c.a*x*x + c.b*y*y + c.c*x*y + c.d*x + c.e*y
}

func (c *ConvexConcane) SolveAnalytical() *Solution {
	sol := &Solution{
		strB: &strings.Builder{},
	}

	fmt.Fprintln(sol.strB, "Kernel function:")
	fmt.Fprintf(sol.strB, kernelFunc, c.a, c.b, c.c, c.d, c.e)

	fmt.Fprintln(sol.strB)

	fmt.Fprintln(sol.strB, "Let us check the feasibility of the conditions for"+
		"the game to belong to the convex-concave class:")
	fmt.Fprintf(sol.strB, "H_xx = 2*a = 2 * %.2f = %.2f < 0\n", c.a, 2*c.a)
	fmt.Fprintf(sol.strB, "H_yy = 2*b = 2 * %.2f = %.2f > 0\n", c.b, 2*c.b)
	fmt.Fprintln(sol.strB, "The game presented is convex-concave.")

	fmt.Fprintln(sol.strB)

	fmt.Fprintf(sol.strB, "H_x = 2ax + cy + d = %.2f * x + %.2f * y + %.2f\n", 2*c.a, c.c, c.d)
	fmt.Fprintf(sol.strB, "H_y = 2by + cx + e = %.2f * y + %.2f * x + %.2f\n", 2*c.b, c.c, c.e)

	fmt.Fprintln(sol.strB)

	x := fmt.Sprintf("%.2fy + %.2f", c.c/-2*c.a, c.d/-2*c.a)
	y := fmt.Sprintf("%.2fx + %.2f", c.c/-2*c.b, c.e/-2*c.b)

	fmt.Fprintf(sol.strB, "x = [cy + d]/[-2a] = %s\n", x)
	fmt.Fprintf(sol.strB, "y = [cx + e]/[-2b] = %s\n", y)

	fmt.Fprintln(sol.strB)

	fmt.Fprintf(sol.strB, "%s, if y >= %.2f\n0 else\n", x, -c.d/c.c)
	fmt.Fprintf(sol.strB, "%s, if x <= %.2f\n0 else\n", y, -c.e/c.c)

	fmt.Fprintln(sol.strB)

	sol.Y = (c.e - (c.c*c.d)/(2*c.a)) / ((c.c*c.c)/(2*c.a) - 2*c.b)
	sol.X = -1 * ((c.c*sol.Y + c.d) / (2 * c.a))
	sol.H = c.h(sol.X, sol.Y)

	return sol
}

func (c *ConvexConcane) SolveNumerical() *Solution {
	sol := &Solution{
		strB: &strings.Builder{},
	}

	N := 2

	lastNIters := make([]float64, 0, lastIters)
	var gameCost, prevGameCost float64
	var xEstimate, yEstimate float64

	for !c.isFinish(lastNIters) {
		m := c.makeMatrix(N)

		fmt.Fprintf(sol.strB, "N = %d\n", N)
		fmt.Fprintln(sol.strB, m.MatrixString())

		highestPrice, highestIdx := m.HighestPrice()
		lowestPrice, lowestIdx := m.LowestPrice()

		if highestPrice != lowestPrice {
			fmt.Fprintln(sol.strB, "Brown-Robinson solution:")

			s := m.SolveBrownRobinson(brownrobinson.Epsilon(0.01))

			xEstimate = float64(slices.Index(s.X, slices.Max(s.X))) / float64(N)
			yEstimate = float64(slices.Index(s.Y, slices.Max(s.Y))) / float64(N)
			gameCost = c.h(xEstimate, yEstimate)
		} else {
			fmt.Fprintln(sol.strB, "Saddle point found:")

			gameCost = highestPrice
			xEstimate = float64(lowestIdx) / float64(N)
			yEstimate = float64(highestIdx) / float64(N)
		}

		fmt.Fprintf(sol.strB, "x = %.3f y = %.3f H = %.3f\n\n",
			xEstimate, yEstimate, gameCost)

		// prevGameCost will appear only since third iteration.
		if N != 2 {
			lastNIters = c.addBuf(lastNIters, math.Abs(gameCost-prevGameCost))
		}

		prevGameCost = gameCost
		N++
	}

	sol.H = gameCost
	sol.X = xEstimate
	sol.Y = yEstimate

	return sol
}

func (c *ConvexConcane) addBuf(last []float64, el float64) []float64 {
	if len(last) == lastIters {
		for i := range len(last) - 1 {
			last[i] = last[i+1]
		}
		last[lastIters-1] = el

		return last
	}

	return append(last, el)
}

func (c *ConvexConcane) isFinish(iters []float64) bool {
	if len(iters) != lastIters {
		return false
	}

	sum := float64(0)

	for _, v := range iters {
		sum += v
	}

	return sum <= eps
}

func (c *ConvexConcane) makeMatrix(N int) *gamematrix.GameMatrix {
	res := make([][]float64, 0, N+1)

	for i := range N + 1 {
		r := make([]float64, 0, N+1)

		for j := range N + 1 {
			r = append(r, c.h(float64(i)/float64(N), float64(j)/float64(N)))
		}

		res = append(res, r)
	}

	m, _ := gamematrix.New(res)

	return m
}
