package main

import (
	"fmt"
	"log"

	"github.com/themilchenko/game_theory/internal/cooperative"
)

const (
	N = 4
)

var charValues = []float64{0, 2, 3, 4, 1, 6, 7, 4, 7, 5, 6, 12, 9, 9, 10, 14}

func main() {
	g, err := cooperative.New(N, charValues)
	if err != nil {
		log.Fatal(err)
	}

	isSuper := g.IsSuperadditiveGame()
	fmt.Println("Игра супераддитивна?", isSuper)

	isConvex, values := g.IsConvex()
	fmt.Println("Игра выпуклая?", isConvex)
	if !isConvex {
		fmt.Printf("Нарушение выпуклости между коалициями %v и %v\n", values[0], values[1])
	}

	shapley, err := g.GetShapleyVector()
	if err != nil {
		log.Fatalf("Ошибка при расчёте вектора Шепли: %w", err)
	}
	fmt.Printf("Вектор Шепли: ( ")
	for _, val := range shapley {
		fmt.Printf("%.3f ", val)
	}
	fmt.Println(")")

	groupRational := g.IsGroupRational()
	fmt.Println("Групповая рационализация выполнена?", groupRational)

	individualRational, vals := g.IsIndividualRational()
	fmt.Println("Индивидуальная рационализация выполнена?", individualRational)
	char := ">="
	if !groupRational {
		char = "<"
	}

	for _, v := range vals {
		fmt.Printf("x_%d=%.3f %s v(%d)=%.3f\n", int(v[0]), v[1], char, int(v[0]), v[2])
	}
}
