package main

import informationalwarfare "github.com/themilchenko/game_theory/internal/informational_warfare"

const (
	N          = 10
	a, b, c, d = 1, 2, 1, 5
	gF, gS     = 1, 2
)

func main() {
	g := informationalwarfare.New(N, a, b, c, d, gF, gS)

	g.Solve(informationalwarfare.General)
}
