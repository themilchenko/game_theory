package analytical

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

type Solution struct {
	x *mat.Dense
	y *mat.Dense
	v float64
}

func (a *Solution) String() string {
	return fmt.Sprintf("x* = %.3v\ny* = %.3v\nv=%.3f",
		mat.Formatted(a.x),
		mat.Formatted(a.y.T()),
		a.v)
}
