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

	dominatorTree := control_flow.NewDominatorTree(cfg)
	expectedImmDom := []uint{0, 0, 0, 0}
	assert.EqualValues(t, expectedImmDom, dominatorTree.ImmDom)
}

func TestBasicDominatorTreeExample(t *testing.T) {
	// Example taken from Henrik Thesis, section 2.6, page 14:
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

	dominatorTree := control_flow.NewDominatorTree(cfg)
	expectedImmDom := []uint{0, 0, 0, 0, 2, 2}
	assert.EqualValues(t, expectedImmDom, dominatorTree.ImmDom)
}
