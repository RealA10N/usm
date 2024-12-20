package control_flow_test

import (
	"testing"

	"alon.kr/x/usm/control_flow"
	"github.com/stretchr/testify/assert"
)

func TestLineDfs(t *testing.T) {
	// 0 -> 1 -> 2

	cfg := control_flow.ControlFlowGraph{
		BasicBlocks: []control_flow.ControlFlowBasicBlock{
			{ForwardEdges: []uint{1}},
			{ForwardEdges: []uint{2}},
			{ForwardEdges: []uint{}},
		},
	}

	result := cfg.Dfs(0)
	assert.EqualValues(t, []uint{0, 1, 2}, result.Preorder)
	assert.EqualValues(t, []uint{0, 0, 1}, result.Parent)
}

func TestBinaryTree(t *testing.T) {
	//    0
	//  1   2
	// 3 4 5 6

	cfg := control_flow.ControlFlowGraph{
		BasicBlocks: []control_flow.ControlFlowBasicBlock{
			{ForwardEdges: []uint{1, 2}}, // 0
			{ForwardEdges: []uint{3, 4}}, // 1
			{ForwardEdges: []uint{5, 6}}, // 2
			{ForwardEdges: []uint{}},     // 3
			{ForwardEdges: []uint{}},     // 4
			{ForwardEdges: []uint{}},     // 5
			{ForwardEdges: []uint{}},     // 6
		},
	}

	result := cfg.Dfs(0)
	//                    index: 0  1  2  3  4  5  6
	assert.EqualValues(t, []uint{0, 1, 4, 2, 3, 5, 6}, result.Preorder)
	assert.EqualValues(t, []uint{0, 0, 0, 1, 1, 2, 2}, result.Parent)
}

func TestSpecialEdges(t *testing.T) {
	//    0      with an additional cross edge 5 -> 1
	//  1   2    a back edge 4 -> 0
	// 3 4 5 6   and a forward edge 0 -> 6

	cfg := control_flow.ControlFlowGraph{
		BasicBlocks: []control_flow.ControlFlowBasicBlock{
			{ForwardEdges: []uint{1, 2, 6}}, // 0
			{ForwardEdges: []uint{3, 4}},    // 1
			{ForwardEdges: []uint{5, 6}},    // 2
			{ForwardEdges: []uint{}},        // 3
			{ForwardEdges: []uint{0}},       // 4
			{ForwardEdges: []uint{1}},       // 5
			{ForwardEdges: []uint{}},        // 6
		},
	}

	result := cfg.Dfs(0)
	//                    index: 0  1  2  3  4  5  6
	assert.EqualValues(t, []uint{0, 1, 4, 2, 3, 5, 6}, result.Preorder)
	assert.EqualValues(t, []uint{0, 0, 0, 1, 1, 2, 2}, result.Parent)
}
