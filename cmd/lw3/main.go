package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/themilchenko/game_theory/internal/bigames"
	"gonum.org/v1/gonum/mat"
)

var PersonalGame = [][]bigames.Position{
	{{X: 10, Y: 7}, {X: 0, Y: 4}},
	{{X: 2, Y: 1}, {X: 9, Y: 3}},
}

func generateRandomGame(n, m int) [][]bigames.Position {
	randGame := make([][]bigames.Position, n)

	for i := range n {
		randGame[i] = make([]bigames.Position, m)

		for j := range m {
			randGame[i][j].X = float64(rand.Intn(101) - 50)
			randGame[i][j].Y = float64(rand.Intn(101) - 50)
		}
	}

	return randGame
}

func main() {
	fmt.Println("Prison Game:")
	fmt.Println(bigames.NewBiGame(bigames.PrisonGame).Solve().String())

	fmt.Println("Family Game:")
	fmt.Println(bigames.NewBiGame(bigames.FamilyGame).Solve().String())

	fmt.Println("Traffic Game:")
	fmt.Println(bigames.NewBiGame(bigames.TrafficGame).Solve().String())

	fmt.Println("Variant Game:")
	fmt.Println(bigames.NewBiGame(PersonalGame).Solve().String())

	s := bigames.NewBiGame(PersonalGame).Solve()
	x, y, v1, v2, err := s.SolveMixedEquilibrium()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(s.String())
	fmt.Printf("X:\n%.3v\nY:\n%.3v\nv1=%.3f\nv2=%.3f\n", mat.Formatted(x), mat.Formatted(y),
		v1, v2)

	fmt.Println("Random Game:")
	fmt.Println(bigames.NewBiGame(generateRandomGame(10, 10)).Solve().String())
}
