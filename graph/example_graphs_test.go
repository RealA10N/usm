package graph_test

import (
	"alon.kr/x/usm/graph"
)

func KnakkegaardGraphExample() graph.Graph {
	// Example taken from Knakkegaard Thesis, section 2.6, page 14:
	// https://users-cs.au.dk/gerth/advising/thesis/henrik-knakkegaard-christensen.pdf

	return graph.NewGraph([][]uint{
		{1, 2}, // 0
		{3},    // 1
		{4, 5}, // 2
		{},     // 3
		{1},    // 4
		{3},    // 5
	})
}

func LengauerTarjanGraphExample() graph.Graph {
	// Example taken from Lengauers's & Tarjan's paper, figures 1 & 2.
	// https://doi.org/10.1145/357062.357071

	R := uint(0)
	A := uint(1)
	B := uint(2)
	C := uint(3)
	D := uint(4)
	E := uint(5)
	F := uint(6)
	G := uint(7)
	H := uint(8)
	I := uint(9)
	J := uint(10)
	K := uint(11)
	L := uint(12)

	return graph.NewGraph([][]uint{
		{A, B, C}, // R
		{D},       // A
		{A, D, E}, // B
		{F, G},    // C
		{L},       // D
		{H},       // E
		{I},       // F
		{I, J},    // G
		{E, K},    // H
		{K},       // I
		{I},       // J
		{R, I},    // K
		{H},       // L
	})
}

func SSABookGraphExample() graph.Graph {
	// Example taken from the SSA Book, figure 3.3(a) & 3.3(b).
	// https://pfalcon.github.io/ssabook/latest/book-full.pdf

	return graph.NewGraph([][]uint{
		{1, 6}, // 0
		{2, 4}, // 1
		{3},    // 2
		{4},    // 3
		{5},    // 4
		{6},    // 5
		{0},    // 6
	})
}

func SreedharGaoGraphExample() graph.Graph {
	// Example taken from Sreedhar's & Gao's paper that first introduced
	// DJ-Graphs (figure 1): https://doi.org/10.1145/199448.199464

	return graph.NewGraph([][]uint{
		{1, 16},     // 0 (START)
		{2, 3, 4},   // 1
		{4, 7},      // 2
		{9},         // 3
		{5},         // 4
		{6},         // 5
		{2, 8},      // 6
		{8},         // 7
		{7, 15},     // 8
		{10, 11},    // 9
		{12},        // 10
		{12},        // 11
		{13},        // 12
		{3, 14, 15}, // 13
		{12},        // 14
		{16},        // 15
		{},          // 16 (END)
	})
}
