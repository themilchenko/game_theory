package bigames

import (
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"gonum.org/v1/gonum/mat"
)

type colorType string

const (
	neshColor    = colorType("\033[31m")
	parettoColor = colorType("\033[33m")
	allColor     = colorType("\033[32m")
	resetColor   = colorType("\033[0m")
)

type Solution struct {
	strB *strings.Builder

	X, Y   *mat.VecDense
	v1, v2 float64

	a, b *mat.Dense
	u    *mat.VecDense
}

func newSolution(m [][]Position) *Solution {
	s := &Solution{
		strB: &strings.Builder{},
	}

	fmt.Fprintf(s.strB, "%sPareto optimal%s\n%sNesh equal%s\n%sAll together%s\n",
		parettoColor, resetColor, neshColor, resetColor, allColor, resetColor)

	t := table.NewWriter()
	t.SetOutputMirror(s.strB)
	t.SetStyle(table.StyleRounded)

	aRaw := make([]float64, 0, len(m)*len(m[0]))
	bRaw := make([]float64, 0, len(m)*len(m[0]))

	for i := range m {
		r := table.Row{}

		for j := range m[i] {
			var color colorType

			if m[i][j].isNesh {
				color = neshColor
			}

			if m[i][j].isPareto {
				color = parettoColor
			}

			if m[i][j].isNesh && m[i][j].isPareto {
				color = allColor
			}

			r = append(r, fmt.Sprintf("%s(%.0f, %.0f)%s",
				color, m[i][j].X, m[i][j].Y, resetColor))

			aRaw = append(aRaw, m[i][j].X)
			bRaw = append(bRaw, m[i][j].Y)
		}

		t.AppendRow(r)
	}

	t.Render()

	u := make([]float64, len(m))
	for i := range u {
		u[i] = 1
	}

	s.a = mat.NewDense(len(m), len(m[0]), aRaw)
	s.b = mat.NewDense(len(m), len(m[0]), bRaw)
	s.u = mat.NewVecDense(len(m), u)

	return s
}

func (s *Solution) String() string {
	fmt.Fprintf(s.strB, "")

	return s.strB.String()
}

func (s *Solution) SolveMixedEquilibrium() (
	x *mat.Dense, y *mat.VecDense, v1, v2 float64, err error,
) {
	rA, cA := s.a.Dims()
	rB, cB := s.b.Dims()

	fmt.Println(mat.Formatted(s.a))
	fmt.Println(mat.Formatted(s.b))

	if rA != cA || rB != cB || rA != rB {
		return nil, nil, 0, 0, fmt.Errorf("матрицы A и B должны быть одинакового размера NxN")
	}

	if s.u.Len() != rA {
		return nil, nil, 0, 0, fmt.Errorf("длина вектора u должна совпадать с размерами матриц A и B")
	}

	var invA, invB mat.Dense

	if err := invA.Inverse(s.a); err != nil {
		return nil, nil, 0, 0, fmt.Errorf("не удалось обратить матрицу A: %v", err)
	}

	if err := invB.Inverse(s.b); err != nil {
		return nil, nil, 0, 0, fmt.Errorf("не удалось обратить матрицу B: %v", err)
	}

	Au := mat.NewVecDense(rA, nil)
	Au.MulVec(&invA, s.u)

	var Bu mat.Dense
	Bu.Mul(&invB, s.u)
	fmt.Println(mat.Formatted(&Bu))

	uTAu := mat.Dot(s.u, Au)

	var uTBu mat.Dense
	uTBu.Mul(s.u.T(), &Bu)
	fmt.Println(mat.Formatted(&uTBu))

	v1 = 1.0 / uTAu
	v2 = 1.0 / uTBu.At(0, 0)

	var uInvB mat.Dense
	uInvB.Mul(s.u.T(), &invB)

	var xVec mat.Dense
	xVec.Scale(v2, &uInvB)

	yVec := mat.NewVecDense(rA, nil)
	yVec.ScaleVec(v1, Au)

	return &xVec, yVec, v1, v2, nil
}
