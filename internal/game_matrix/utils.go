package gamematrix

import (
	"slices"

	"gonum.org/v1/gonum/mat"
)

func covertMatVec(v mat.Vector) []float64 {
	converted := make([]float64, 0, v.Len())

	for i := range v.Len() {
		converted = append(converted, v.AtVec(i))
	}

	return converted
}

func findMinInVec(v mat.Vector) float64 {
	return slices.Min(covertMatVec(v))
}

func findMaxInVec(v mat.Vector) float64 {
	return slices.Max(covertMatVec(v))
}
