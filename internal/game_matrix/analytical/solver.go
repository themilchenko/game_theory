package analytical

import (
	"gonum.org/v1/gonum/mat"
)

type Solver struct {
	m           *mat.Dense
	mInv        *mat.Dense
	u           *mat.VecDense
	denominator *mat.Dense
}

func New(m *mat.Dense) (*Solver, error) {
	a := &Solver{
		m: m,
	}

	a.mInv = &mat.Dense{}
	if err := a.mInv.Inverse(a.m); err != nil {
		return nil, err
	}

	r := a.m.RawMatrix().Rows

	oneVec := make([]float64, 0, r)
	for range r {
		oneVec = append(oneVec, 1)
	}

	a.u = mat.NewVecDense(r, oneVec)

	var uTCInv mat.Dense
	uTCInv.Mul(a.u.T(), a.mInv)

	a.denominator = &mat.Dense{}
	a.denominator.Mul(&uTCInv, a.u)

	return a, nil
}

func (a *Solver) Solve() (*Solution, error) {
	x, err := a.solveX()
	if err != nil {
		return nil, err
	}

	y, err := a.solveY()
	if err != nil {
		return nil, err
	}

	v, err := a.solveV()
	if err != nil {
		return nil, err
	}

	return &Solution{
		x: x,
		y: y,
		v: v,
	}, nil
}

func (a *Solver) solveX() (*mat.Dense, error) {
	var numerator mat.Dense
	numerator.Mul(a.u.T(), a.mInv)

	var res mat.Dense
	res.Scale(1/a.denominator.At(0, 0), &numerator)

	return &res, nil
}

func (s *Solver) solveY() (*mat.Dense, error) {
	var numerator mat.Dense
	numerator.Mul(s.mInv, s.u)

	var res mat.Dense
	res.Scale(1/s.denominator.At(0, 0), &numerator)

	return &res, nil
}

func (s *Solver) solveV() (float64, error) {
	return 1 / s.denominator.At(0, 0), nil
}
