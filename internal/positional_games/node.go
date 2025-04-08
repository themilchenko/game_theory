package positionalgames

type node struct {
	id int

	v        []tuple
	parent   *node
	children []*node

	playerNum int

	isBestChildren []*node
}

type tuple []int
