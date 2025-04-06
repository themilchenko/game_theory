package bigames

type Position struct {
	X float64
	Y float64

	isPareto bool
	isNesh   bool
}

type BiGame struct {
	m [][]Position
}

func NewBiGame(m [][]Position) *BiGame {
	g := &BiGame{
		m: make([][]Position, len(m)),
	}

	for i := range m {
		g.m[i] = make([]Position, 0, len(m[i]))

		g.m[i] = append(g.m[i], m[i]...)
	}

	return g
}

func (g *BiGame) Solve() *Solution {
	for i, s := range g.m {
		for j := range s {
			if g.isParetoOptimal(i, j) {
				g.m[i][j].isPareto = true
			}

			if g.isNeshEqual(i, j) {
				g.m[i][j].isNesh = true
			}
		}
	}

	return newSolution(g.m)
}

func (g *BiGame) isParetoOptimal(a, b int) bool {
	for i := range g.m {
		for j := range g.m[i] {
			if g.notParetoCondition(i, j, a, b) {
				return false
			}
		}
	}

	return true
}

func (g *BiGame) notParetoCondition(i, j, a, b int) bool {
	return (g.m[i][j].X >= g.m[a][b].X && g.m[i][j].Y > g.m[a][b].Y) ||
		(g.m[i][j].X > g.m[a][b].X && g.m[i][j].Y >= g.m[a][b].Y)
}

func (g *BiGame) isNeshEqual(a, b int) bool {
	for i := range g.m {
		if g.m[i][b].X >= g.m[a][b].X && i != a {
			return false
		}
	}

	for j := range g.m[a] {
		if g.m[a][j].Y >= g.m[a][b].Y && j != b {
			return false
		}
	}

	return true
}
