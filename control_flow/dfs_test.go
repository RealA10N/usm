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

	preorder := cfg.PreOrderDfs()
	assert.Len(t, preorder, len(cfg.BasicBlocks))
	assert.EqualValues(t, []uint{0, 1, 2}, preorder)
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

	preorder := cfg.PreOrderDfs()
	assert.Len(t, preorder, len(cfg.BasicBlocks))
	//                    index: 0  1  2  3  4  5  6
	assert.EqualValues(t, []uint{0, 1, 4, 2, 3, 5, 6}, preorder)
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

	preorder := cfg.PreOrderDfs()
	assert.Len(t, preorder, len(cfg.BasicBlocks))
	//                    index: 0  1  2  3  4  5  6
	assert.EqualValues(t, []uint{0, 1, 4, 2, 3, 5, 6}, preorder)

}
