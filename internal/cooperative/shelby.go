package cooperative

import (
	"errors"
	"fmt"
	"math"
)

type CooperativeGame struct {
	players        int
	totalCoalition []int
	charMapping    map[string]float64
	shapleyVec     []float64
	allCoalitions  [][]int
	coalitionPairs [][2][]int
}

func New(n int, charValues []float64) (*CooperativeGame, error) {
	if len(charValues) != 1<<n {
		return nil, fmt.Errorf("length of charValues should be %d, got %d", 1<<n, len(charValues))
	}

	g := &CooperativeGame{
		players:        n,
		totalCoalition: make([]int, n),
		charMapping:    make(map[string]float64),
	}

	for i := range g.totalCoalition {
		g.totalCoalition[i] = i + 1
	}

	g.generateCoalitions()

	for i, c := range g.allCoalitions {
		key := coalitionKey(c)
		g.charMapping[key] = charValues[i]
	}

	g.generateCoalitionPairs()

	return g, nil
}

func (g *CooperativeGame) GetShapleyVector() ([]float64, error) {
	if !g.IsSuperadditiveGame() {
		return nil, errors.New("game is not superadditive")
	}
	if g.shapleyVec != nil {
		return g.shapleyVec, nil
	}

	n := g.players
	g.shapleyVec = make([]float64, n)
	totalPerm := factorial(n)

	for i := 1; i <= n; i++ {
		var sum float64

		for _, S := range g.allCoalitions {
			if contains(S, i) {
				Swithout := remove(S, i)
				weight := float64(factorial(len(S)-1)) * float64(factorial(n-len(S)))
				vDiff := float64(g.charFunction(S)) - float64(g.charFunction(Swithout))
				sum += float64(weight * vDiff)
			}
		}
		g.shapleyVec[i-1] = float64(float64(1)/float64(totalPerm)) * sum
	}
	return g.shapleyVec, nil
}

func (g *CooperativeGame) IsGroupRational() bool {
	sum := 0.0
	for _, v := range g.shapleyVec {
		sum += float64(v)
	}
	totalValue := float64(g.charFunction(g.totalCoalition))

	return math.Abs(sum-totalValue) < 1e-9
}

func (g *CooperativeGame) IsIndividualRational() (bool, [][3]float64) {
	res := make([][3]float64, 0)

	for i := 1; i <= g.players; i++ {
		if g.shapleyVec[i-1] < g.charFunction([]int{i}) {
			return false, [][3]float64{{float64(i), g.shapleyVec[i-1], g.charFunction([]int{i})}}
		}

		res = append(res, [3]float64{float64(i), g.shapleyVec[i-1], g.charFunction([]int{i})})
	}
	return true, res
}

func (g *CooperativeGame) charFunction(c []int) float64 {
	return g.charMapping[coalitionKey(c)]
}

func (g *CooperativeGame) IsSuperadditiveGame() bool {
	for _, pair := range g.coalitionPairs {
		A := toSet(pair[0])
		B := toSet(pair[1])
		if disjoint(A, B) {
			union := toSlice(unionSet(A, B))
			vUnion := g.charFunction(union)
			vA := g.charFunction(pair[0])
			vB := g.charFunction(pair[1])
			if vUnion < vA+vB {
				return false
			}
		}
	}
	return true
}

func (g *CooperativeGame) IsConvex() (bool, [2][]int) {
	for _, pair := range g.coalitionPairs {
		A := toSet(pair[0])
		B := toSet(pair[1])
		union := toSlice(unionSet(A, B))
		inter := toSlice(intersectionSet(A, B))
		vUnion := g.charFunction(union)
		vInter := g.charFunction(inter)
		vA := g.charFunction(pair[0])
		vB := g.charFunction(pair[1])

		if vUnion+vInter < vA+vB {
			return false, pair
		}
	}
	return true, [2][]int{}
}

func (g *CooperativeGame) generateCoalitions() {
	g.allCoalitions = [][]int{}
	players := g.totalCoalition
	n := len(players)

	for k := 0; k <= n; k++ {
		g.combine(players, k, []int{}, 0)
	}
}

func (g *CooperativeGame) combine(players []int, k int, curr []int, start int) {
	if len(curr) == k {
		coal := append([]int{}, curr...)
		g.allCoalitions = append(g.allCoalitions, coal)
		return
	}
	for i := start; i < len(players); i++ {
		g.combine(players, k, append(curr, players[i]), i+1)
	}
}

func (g *CooperativeGame) generateCoalitionPairs() {
	for i := range len(g.allCoalitions) {
		for j := i + 1; j < len(g.allCoalitions); j++ {
			g.coalitionPairs = append(g.coalitionPairs, [2][]int{g.allCoalitions[i], g.allCoalitions[j]})
		}
	}
}
