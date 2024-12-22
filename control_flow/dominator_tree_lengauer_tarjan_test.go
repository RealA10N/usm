package control_flow_test

import (
	"testing"

	"alon.kr/x/usm/control_flow"
	"github.com/stretchr/testify/assert"
)

func TestIfElseDominatorTreeBuilder(t *testing.T) {
	// control flow graph:
	//     0
	//    / \
	//   1   2
	//    \ /
	//     3
	//
	// dominator tree:
	//      0
	//    / | \
	//   1  2  3

	cfg := control_flow.NewGraph([][]uint{{1, 2}, {3}, {3}, {}})
	dominatorTree := cfg.DominatorTree(0)

	expectedImmDom := []uint{0, 0, 0, 0}
	assert.EqualValues(t, expectedImmDom, dominatorTree.ImmDom)
}

func TestKnakkegaardsDominatorTreeExample(t *testing.T) {
	// Example taken from Knakkegaard Thesis, section 2.6, page 14:
	// https://users-cs.au.dk/gerth/advising/thesis/henrik-knakkegaard-christensen.pdf

	cfg := KnakkegaardGraphExample()
	dominatorTree := cfg.DominatorTree(0)

	expectedImmDom := []uint{0, 0, 0, 0, 2, 2}
	assert.EqualValues(t, expectedImmDom, dominatorTree.ImmDom)
}

func TestLengauerTarjansPaperDominatorTreeExample(t *testing.T) {
	// Example taken from Lengauers's & Tarjan's paper, figures 1 & 2.
	// https://doi.org/10.1145/357062.357071

	cfg := LengauerTarjanGraphExample()
	dominatorTree := cfg.DominatorTree(0)
	R := uint(0)
	C := uint(3)
	D := uint(4)
	G := uint(7)

	//                       R  A  B  C  D  E  F  G  H  I  J  K  L
	expectedImmDom := []uint{R, R, R, R, R, R, C, C, R, R, G, R, D}
	assert.EqualValues(t, expectedImmDom, dominatorTree.ImmDom)
}

func TestSsaBookDominatorTreeExample(t *testing.T) {
	// Example taken from the SSA Book, figure 3.3(a) & 3.3(b).
	// https://pfalcon.github.io/ssabook/latest/book-full.pdf
	//
	// control flow graph:
	//        0
	//        | \
	//        |   1
	//        |   | \
	//        |   |   2
	//        |   |   |
	//        |   |   3
	//        |   | /
	//        |   4
	//        |   |
	//        |   5
	//        | /
	//        6
	//
	// dominator tree:
	//        0
	//       / \
	//      6   1
	//         / \
	//        4   2
	//        |   |
	//        5   3

	cfg := SSABookGraphExample()
	dominatorTree := cfg.DominatorTree(0)

	//                       0  1  2  3  4  5  6
	expectedImmDom := []uint{0, 0, 1, 2, 1, 4, 0}
	assert.EqualValues(t, expectedImmDom, dominatorTree.ImmDom)
}

func TestDjGraphPaperDominatorTreeExample(t *testing.T) {
	// Example taken from Sreedhar's & Gao's paper that first introduced
	// DJ-Graphs (figure 1): https://doi.org/10.1145/199448.199464

	cfg := SreedharGaoGraphExample()
	dominatorTree := cfg.DominatorTree(0)
	//                       0  1  2  3  4  5  6  7  8  9 10 11 12  13  14 15 16
	expectedImmDom := []uint{0, 0, 1, 1, 1, 4, 5, 1, 1, 3, 9, 9, 9, 12, 13, 1, 0}
	assert.EqualValues(t, expectedImmDom, dominatorTree.ImmDom)
}
