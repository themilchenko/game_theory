package main

import (
	"fmt"
	"log"

	convexconcane "github.com/themilchenko/game_theory/internal/convex-concane"
)

var (
	a float64 = -10
	b float64 = 40 / float64(3)
	c float64 = 40
	d float64 = -16
	e float64 = -32
)

func main() {
	c, err := convexconcane.New(a, b, c, d, e)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(c.SolveAnalytical())
	fmt.Println(c.SolveNumerical())
}
