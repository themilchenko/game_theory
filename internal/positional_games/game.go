package positionalgames

import (
	"errors"
	"math/rand"
	"slices"
)

type Game struct {
	r         *node
	curPlayer int
	playerNum int
	leafLvl   []*node
	bestStrat [][]*node
}

func NewGame(depth, playerNum int, strategyNum []int, winRange [2]int) (*Game, error) {
	if len(strategyNum) != playerNum {
		return nil, errors.New("length of startegies number should be the same as palyers number")
	}

	idCounter := 1

	g := &Game{
		r: &node{
			id: idCounter,
		},
		playerNum: playerNum,
		bestStrat: make([][]*node, 0),
	}

	idCounter++

	curLvl := []*node{g.r}

	for i := range depth {
		g.curPlayer = i % playerNum
		childrenLen := strategyNum[i%len(strategyNum)]

		newLvl := make([]*node, 0, len(curLvl))

		for _, child := range curLvl {
			child.children = make([]*node, childrenLen)
			child.playerNum = g.curPlayer

			for k := range child.children {
				child.children[k] = &node{
					id:     idCounter,
					parent: child,
				}

				idCounter++
			}

			newLvl = append(newLvl, child.children...)
		}

		curLvl = newLvl
	}

	for _, node := range curLvl {
		node.v = make([]tuple, 1)
		node.v[0] = make([]int, playerNum)

		for i := range node.v[0] {
			node.v[0][i] = rand.Intn(winRange[1]-winRange[0]+1) + winRange[0]
		}
	}

	g.leafLvl = curLvl

	return g, nil
}

func (g *Game) Solve() []tuple {
	curLvl := g.getLvl(g.leafLvl)

	for {
		for _, n := range curLvl {
			res := g.compare(n.children, g.curPlayer)
			for _, v := range res {
				n.v = append(n.v, v.v...)
			}
		}

		if len(curLvl) == 1 {
			break
		}

		g.curPlayer = (g.curPlayer + 1) % g.playerNum
		curLvl = g.getLvl(curLvl)
	}

	return g.r.v
}

func (g *Game) getLvl(children []*node) []*node {
	lvl := make([]*node, 0)

	for _, c := range children {
		if !slices.Contains(lvl, c.parent) {
			lvl = append(lvl, c.parent)
		}
	}

	return lvl
}

func (g *Game) compare(nodes []*node, player int) []*node {
	type val struct {
		v    tuple
		node *node
	}

	maxS := make([]val, 0, len(nodes))

	// From every node extract max tuple.
	for _, n := range nodes {
		maxS = append(maxS, val{
			v: slices.MaxFunc(n.v, func(a, b tuple) int {
				return a[player] - b[player]
			}),
			node: n,
		})
	}

	// Check for max tuple.
	m := slices.MaxFunc(maxS, func(a, b val) int {
		return a.v[player] - b.v[player]
	})

	res := make([]*node, 0)

	// Check that there is no tuple with max value.
	for _, v := range maxS {
		if v.v[player] == m.v[player] {
			res = append(res, v.node)
		}
	}

	return res
}
