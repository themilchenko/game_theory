package main

import (
	"fmt"
	"log"
	"os"

	positionalgames "github.com/themilchenko/game_theory/internal/positional_games"
)

var (
	depth       = 7
	players     = 2
	strategyNum = []int{2, 3}
	winRange    = [2]int{-5, 25}
)

func main() {
	g, err := positionalgames.NewGame(depth, players, strategyNum, winRange)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(g.Solve())

	f, err := os.Create("./artifacts/lw4/tree.dot")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fmt.Fprintln(f, g.Dot())
}
