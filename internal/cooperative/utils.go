package cooperative

import (
	"fmt"
	"sort"
)

func factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * factorial(n-1)
}

func coalitionKey(c []int) string {
	cCopy := append([]int(nil), c...)
	sort.Ints(cCopy)
	key := ""
	for _, i := range cCopy {
		key += fmt.Sprintf("%d,", i)
	}
	return key
}

func toSet(c []int) map[int]bool {
	m := make(map[int]bool)
	for _, v := range c {
		m[v] = true
	}
	return m
}

func toSlice(m map[int]bool) []int {
	s := []int{}
	for k := range m {
		s = append(s, k)
	}
	sort.Ints(s)
	return s
}

func disjoint(a, b map[int]bool) bool {
	for k := range a {
		if b[k] {
			return false
		}
	}
	return true
}

func unionSet(a, b map[int]bool) map[int]bool {
	res := make(map[int]bool)
	for k := range a {
		res[k] = true
	}
	for k := range b {
		res[k] = true
	}
	return res
}

func intersectionSet(a, b map[int]bool) map[int]bool {
	res := make(map[int]bool)
	for k := range a {
		if b[k] {
			res[k] = true
		}
	}
	return res
}

func contains(c []int, x int) bool {
	for _, v := range c {
		if v == x {
			return true
		}
	}
	return false
}

func remove(c []int, x int) []int {
	newC := []int{}
	for _, v := range c {
		if v != x {
			newC = append(newC, v)
		}
	}
	return newC
}
