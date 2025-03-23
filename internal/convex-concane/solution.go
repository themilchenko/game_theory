package convexconcane

import (
	"fmt"
	"strings"
)

type Solution struct {
	X float64
	Y float64
	H float64

	strB *strings.Builder
}

func (s *Solution) String() string {
	fmt.Fprintf(s.strB, "x = %.2f\ny = %.2f\nH = %.2f", s.Y, s.X, s.H)

	return s.strB.String()
}
