package informationalwarfare

import (
	"fmt"
	"math"
	"math/rand"
	"sort"

	"gonum.org/v1/gonum/mat"
)

type SolutionType int

const (
	General SolutionType = iota
	TargetFunction

	absMaxInfuence = 100
	diff           = 20
)

type InformationalWarfare struct {
	n, a, b, c, d, gF, gS int
}

func New(n, a, b, c, d, gF, gS int) *InformationalWarfare {
	return &InformationalWarfare{
		n:  n,
		a:  a,
		b:  b,
		c:  c,
		d:  d,
		gF: gF,
		gS: gS,
	}
}

func (g *InformationalWarfare) Solve(solType SolutionType) {
	m := g.generateTrustMatrix()

	fmt.Printf("%.3v\n", mat.Formatted(m))

	if !g.isStochastic(m) {
		fmt.Println("matrix is not stochastic")
	}

	a, b := 1, 20
	initialOpinions := g.generateInitialOpinions(a, b+1, g.n)

	fmt.Printf("Random vector of agent's opinions: %v\n", mat.Formatted(initialOpinions.T()))

	finalOpinions, resTrustMatrix := g.computeFinalOpinions(m, initialOpinions, 0.000001)

	fmt.Printf("Final opinions:\n %.3v\nRes trust matrix:\n %.3v\n",
		mat.Formatted(finalOpinions.T()), mat.Formatted(resTrustMatrix))

	player1InfluenceIndices, player2InfluenceIndices := g.getPlayerAgentIndices(g.n)

	fmt.Printf("Indices of agents of influence for the first player: %v\n", player1InfluenceIndices)
	fmt.Printf("Indices of agents of influence for the second player: %v\n", player2InfluenceIndices)

	player1Influence := rand.Intn(absMaxInfuence)
	player2Influence := -rand.Intn(absMaxInfuence)

	fmt.Println("Initial opinion of 1 player's agents: ", player1Influence)
	fmt.Println("Initial opinion of 2 player's agents: ", player2Influence)

	for _, idx := range player1InfluenceIndices {
		initialOpinions.SetVec(idx, float64(player1Influence))
	}

	for _, idx := range player2InfluenceIndices {
		initialOpinions.SetVec(idx, float64(player2Influence))
	}

	finalOpinions, resTrustMatrix = g.computeFinalOpinions(m, initialOpinions, 0.000001)
	fmt.Printf("After the agents interact, the opinion vector converges to the value of X =\n%.3v\n", mat.Formatted(finalOpinions.T()))
	fmt.Printf("Result trust matrix:\n%.3v\n", mat.Formatted(resTrustMatrix))

	switch solType {
	case General:
		g.WithGeneral(finalOpinions, player1Influence, player2Influence)
	case TargetFunction:
		g.WithTargetFunction(player1InfluenceIndices, player2InfluenceIndices, resTrustMatrix)
	default:
		fmt.Println("Undefined solution type")
	}
}

func (g *InformationalWarfare) WithGeneral(finalOpinions *mat.VecDense, player1Influence, player2Influence int) {
	diff1 := math.Abs(finalOpinions.AtVec(0) - float64(player1Influence))
	diff2 := math.Abs(finalOpinions.AtVec(0) - float64(player2Influence))

	if diff1 < diff2 {
		fmt.Println("First player won")
	} else {
		fmt.Println("Second player won")
	}
}

func (g *InformationalWarfare) WithTargetFunction(player1InfluenceIndices, player2InfluenceIndices []int,
	resTrustMatrix *mat.Dense,
) {
	fmt.Printf("Result trust matrix:\n%.3v\n", mat.Formatted(resTrustMatrix))

	rF := 0.0
	for _, idx := range player1InfluenceIndices {
		rF += resTrustMatrix.At(0, idx)
	}

	rS := 0.0
	for _, idx := range player2InfluenceIndices {
		rS += resTrustMatrix.At(1, idx)
	}

	fmt.Printf("(r_f, r_s) = (%.3v, %.3v)\n", rF, rS)

	targetF := fmt.Sprintf("%.3f * u^2 + %.3f * u + %.3f * v + %.3f * u * v + %.3f v^2",
		-float64(g.b)*rF*rF-float64(g.gF)/2,
		float64(g.a)*rF,
		float64(g.a)*rS,
		-float64(2)*float64(g.b)*rF*rS,
		-float64(g.b)*rS*rS,
	)

	targetS := fmt.Sprintf("%.3f * u^2 + %.3f * u + %.3f * v + %.3f * u * v + %.3f v^2",
		-float64(g.d)*rF*rF,
		float64(g.c)*rF,
		float64(g.c)*rS,
		-float64(2)*float64(g.d)*rF*rS,
		-float64(g.d)*rS*rS-(float64(g.gS)/2),
	)

	fmt.Printf("Фf(u, v) = %s\n", targetF)
	fmt.Printf("Фs(u, v) = %s\n", targetS)

	targetFDiffU := fmt.Sprintf("%.3f * u + %.3f * v + %.3f",
		-2*float64(g.b)*rF*rF-float64(g.gF),
		-2*float64(g.b)*rS*rF,
		float64(g.a)*rF,
	)
	targetSDiffV := fmt.Sprintf("%.3f * v + %.3f * u + %.3f",
		-2*float64(g.d)*rS*rS-float64(g.gS),
		-float64(2)*float64(g.d)*rF*rS,
		float64(g.c)*rS,
	)

	fmt.Printf("Фf(u, v)|'_u = %s\n", targetFDiffU)
	fmt.Printf("Фs(u, v)|'_v = %s\n", targetSDiffV)

	u := (2*rS*rS*float64(g.d)*float64(g.a)*rF + float64(g.a)*rF*float64(g.gS) - 2*rF*rS*rS*float64(g.c)*float64(g.b)) /
		(float64(g.gF)*float64(g.gS) + 2*float64(g.b)*rF*rF*float64(g.gS) + 2*float64(g.d)*rS*rS*float64(g.gF))
	v := (u*(-2*float64(g.b)*rF*rF-float64(g.gF)) + float64(g.a)*rF) / (2 * float64(g.b) * rF * rS)

	fmt.Printf("v = %.3f, u = %.3f\n", v, u)

	X := u*rF + v*rS
	fmt.Printf("X = %.3f\n", X)

	xMaxF, xMaxS := float64(g.a)/float64((2*g.b)), float64(g.c)/float64((2*g.d))
	deltaXF := math.Abs(X - float64(xMaxF))
	deltaXS := math.Abs(X - float64(xMaxS))

	fmt.Printf("X_max_f = %.3f, X_max_s = %.3f\n Let's find distance:\ndelta_x_f = %.3f, delta_x_s = %.3f\n",
		xMaxF, xMaxS, deltaXF, deltaXS)

	if deltaXF < deltaXS {
		fmt.Println("First player won")
	} else {
		fmt.Println("Second player won")
	}
}

func (g *InformationalWarfare) generateTrustMatrix() *mat.Dense {
	data := make([]float64, g.n*g.n)

	for i := range g.n {
		var rowSum float64
		for j := range g.n {
			val := rand.Float64()
			data[i*g.n+j] = val
			rowSum += val
		}

		for j := range g.n {
			data[i*g.n+j] /= rowSum
		}
	}

	return mat.NewDense(g.n, g.n, data)
}

func (g *InformationalWarfare) isStochastic(matrix *mat.Dense) bool {
	for i := range g.n {
		row := matrix.RawRowView(i)
		var rowSum float64
		for _, val := range row {
			rowSum += val
		}

		if !(rowSum > 0.999999 && rowSum < 1.000001) {
			return false
		}
	}
	return true
}

func (g *InformationalWarfare) generateInitialOpinions(minVal, maxVal, N int) *mat.VecDense {
	data := make([]float64, N)
	for i := range N {
		data[i] = float64(rand.Intn(maxVal-minVal+1) + minVal)
	}
	return mat.NewVecDense(N, data)
}

func (g *InformationalWarfare) computeFinalOpinions(trustMatrix *mat.Dense,
	initialOpinions *mat.VecDense, epsilon float64,
) (*mat.VecDense, *mat.Dense) {
	n, _ := trustMatrix.Dims()

	currentOpinions := mat.NewVecDense(n, nil)
	currentOpinions.CopyVec(initialOpinions)
	previousOpinions := mat.NewVecDense(n, nil)

	resTrustMatrix := mat.DenseCopyOf(trustMatrix)

	i := 0

	for {
		difference := mat.NewVecDense(n, nil)
		difference.SubVec(currentOpinions, previousOpinions)
		norm := mat.Norm(difference, 2)
		if norm <= epsilon {
			break
		}

		previousOpinions.CopyVec(currentOpinions)

		currentOpinions.MulVec(trustMatrix, currentOpinions)

		temp := mat.NewDense(n, n, nil)
		temp.Mul(resTrustMatrix, trustMatrix)
		resTrustMatrix = temp

		fmt.Printf("x(%d)=%.4v\n", i, mat.Formatted(currentOpinions.T()))

		i++
	}

	return currentOpinions, resTrustMatrix
}

func (g *InformationalWarfare) getPlayerAgentIndices(nAgents int) ([]int, []int) {
	if nAgents < 2 {
		return []int{}, []int{}
	}
	player1Count := rand.Intn(nAgents-1) + 1
	player2Count := rand.Intn(nAgents-player1Count) + 1

	perm := rand.Perm(nAgents)
	player1Indices := make([]int, player1Count)
	player2Indices := make([]int, player2Count)
	copy(player1Indices, perm[:player1Count])
	copy(player2Indices, perm[player1Count:player1Count+player2Count])

	sort.Ints(player1Indices)
	sort.Ints(player2Indices)
	return player1Indices, player2Indices
}
