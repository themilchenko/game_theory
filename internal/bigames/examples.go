package bigames

var (
	e = float64(0.001)

	PrisonGame = [][]Position{
		{{X: -5, Y: -5}, {X: 0, Y: -10}},
		{{X: -10, Y: 0}, {X: -1, Y: -1}},
	}

	FamilyGame = [][]Position{
		{{X: 4, Y: 1}, {X: 0, Y: 0}},
		{{X: 0, Y: 0}, {X: 1, Y: 4}},
	}

	TrafficGame = [][]Position{
		{{X: 1, Y: 1}, {X: 1 - e, Y: 2}},
		{{X: 2, Y: 1 - e}, {X: 0, Y: 0}},
	}
)
