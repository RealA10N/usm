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
	//
	//      0
	//    / | \
	//   1  2  3

	cfg := control_flow.ControlFlowGraph{
		BasicBlocks: []control_flow.ControlFlowBasicBlock{
			{ForwardEdges: []uint{1, 2}, BackwardEdges: []uint{}}, // 0
			{ForwardEdges: []uint{3}, BackwardEdges: []uint{0}},   // 1
			{ForwardEdges: []uint{3}, BackwardEdges: []uint{0}},   // 2
			{ForwardEdges: []uint{}, BackwardEdges: []uint{1, 2}}, // 3
		},
	}

	dominatorTree := cfg.DominatorTree()
	expectedImmDom := []uint{0, 0, 0, 0}
	assert.EqualValues(t, expectedImmDom, dominatorTree.ImmDom)
}

func TestKnakkegaardsDominatorTreeExample(t *testing.T) {
	// Example taken from Knakkegaard Thesis, section 2.6, page 14:
	// https://users-cs.au.dk/gerth/advising/thesis/henrik-knakkegaard-christensen.pdf
	//
	// control flow graph:
	//       0
	//      ↙ ↘
	//     1   2
	//    ↙ ↖ ↙ ↘
	//   3   4   5
	//    ↖←←←←←↙
	//
	// dominator tree:
	//       0
	//     / | \
	//    1  2  3
	//      / \
	//     4   5

	cfg := control_flow.ControlFlowGraph{
		BasicBlocks: []control_flow.ControlFlowBasicBlock{
			{ForwardEdges: []uint{1, 2}, BackwardEdges: []uint{}},  // 0
			{ForwardEdges: []uint{3}, BackwardEdges: []uint{0, 4}}, // 1
			{ForwardEdges: []uint{4, 5}, BackwardEdges: []uint{0}}, // 2
			{ForwardEdges: []uint{}, BackwardEdges: []uint{1, 5}},  // 3
			{ForwardEdges: []uint{1}, BackwardEdges: []uint{2}},    // 4
			{ForwardEdges: []uint{3}, BackwardEdges: []uint{2}},    // 5
		},
	}

	dominatorTree := cfg.DominatorTree()
	expectedImmDom := []uint{0, 0, 0, 0, 2, 2}
	assert.EqualValues(t, expectedImmDom, dominatorTree.ImmDom)
}

func TestTarjansDominatorTreeExample(t *testing.T) {
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

	cfg := control_flow.ControlFlowGraph{
		BasicBlocks: []control_flow.ControlFlowBasicBlock{
			{ForwardEdges: []uint{A, B, C}, BackwardEdges: []uint{K}},    // R
			{ForwardEdges: []uint{D}, BackwardEdges: []uint{R, B}},       // A
			{ForwardEdges: []uint{A, D, E}, BackwardEdges: []uint{R}},    // B
			{ForwardEdges: []uint{F, G}, BackwardEdges: []uint{R}},       // C
			{ForwardEdges: []uint{L}, BackwardEdges: []uint{A, B}},       // D
			{ForwardEdges: []uint{H}, BackwardEdges: []uint{H, B}},       // E
			{ForwardEdges: []uint{I}, BackwardEdges: []uint{C}},          // F
			{ForwardEdges: []uint{I, J}, BackwardEdges: []uint{C}},       // G
			{ForwardEdges: []uint{E, K}, BackwardEdges: []uint{L, E}},    // H
			{ForwardEdges: []uint{K}, BackwardEdges: []uint{K, F, G, J}}, // I
			{ForwardEdges: []uint{I}, BackwardEdges: []uint{G}},          // J
			{ForwardEdges: []uint{R, I}, BackwardEdges: []uint{H, I}},    // K
			{ForwardEdges: []uint{H}, BackwardEdges: []uint{D}},          // L
		},
	}

	dominatorTree := cfg.DominatorTree()
	//                       R  A  B  C  D  E  F  G  H  I  J  K  L
	expectedImmDom := []uint{R, R, R, R, R, R, C, C, R, R, G, R, D}
	assert.EqualValues(t, expectedImmDom, dominatorTree.ImmDom)
}
