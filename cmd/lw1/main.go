package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	gamematrix "github.com/themilchenko/game_theory/internal/game_matrix"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: ./path/to/exec task.json")
	}

	if _, err := os.Stat(os.Args[1]); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Fatal(fmt.Errorf("file %q not exists: %w", os.Args[1], err))
		}
		log.Fatal(err)
	}

	fContent, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(fmt.Errorf("failed to read %q file: %w", os.Args[1], err))
	}

	var matrix [][]float64

	if err := json.Unmarshal(fContent, &matrix); err != nil {
		log.Fatal(fmt.Errorf("failed to parse matrix from json: %w", err))
	}

	game, err := gamematrix.New(matrix)
	if err != nil {
		log.Fatal(fmt.Errorf("can't creage game matrix: %w", err))
	}

	fmt.Println(game.String())
	fmt.Printf("Lowest Price: %f\n", game.LowestPrice())
	fmt.Printf("Highest Price: %f\n", game.HighestPrice())

	fmt.Println()

	sol, err := game.SolveAnalytical()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to solve analytical: %w", err))
	}

	fmt.Println("Analytical solution:")
	fmt.Println(sol.String())

	fmt.Println("Brown Robinson solution:")

	fmt.Println(game.SolveBrownRobinson().String())
}
