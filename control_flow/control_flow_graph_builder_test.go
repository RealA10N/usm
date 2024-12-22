package control_flow_test

import (
	"testing"

	"alon.kr/x/usm/control_flow"
	"github.com/stretchr/testify/assert"
)

func TestBuildSingleton(t *testing.T) {
	g := control_flow.NewGraph(1, [][]uint{
		{},
	})
	cfg := g.ControlFlowGraph(0)
	assert.EqualValues(t, 1, cfg.Graph.Size())
	assert.Empty(t, cfg.Nodes[0].ForwardEdges)
	assert.Empty(t, cfg.Nodes[0].BackwardEdges)
}

func TestBuildSelfLoop(t *testing.T) {
	g := control_flow.NewGraph(1, [][]uint{
		{0},
	})

	cfg := g.ControlFlowGraph(0)
	assert.EqualValues(t, 1, cfg.Graph.Size())
	assert.EqualValues(t, []uint{0}, cfg.Nodes[0].ForwardEdges)
	assert.EqualValues(t, []uint{0}, cfg.Nodes[0].BackwardEdges)
}

func TestSingleBlock(t *testing.T) {
	// input graph:
	// 0 -> 1 -> 2
	//
	// control flow graph:
	// 0 (single node, no edges)

	g := control_flow.NewGraph(3, [][]uint{
		{1},
		{2},
		{},
	})

	cfg := g.ControlFlowGraph(0)
	assert.EqualValues(t, 1, cfg.Size())
	assert.EqualValues(t, []uint{0, 1, 2}, cfg.BasicBlockToNodes[0])
	assert.EqualValues(t, []uint{0, 0, 0}, cfg.NodeToBasicBlock)
	assert.Empty(t, cfg.Nodes[0].ForwardEdges)
	assert.Empty(t, cfg.Nodes[0].BackwardEdges)
}

func TestSimpleLoop(t *testing.T) {
	// input graph:
	// 0 -> 1 -> 2
	//      ^----+
	//
	// control flow graph:
	// 0 -> 1-+
	//      ^-+

	g := control_flow.NewGraph(3, [][]uint{
		{1},
		{2},
		{1},
	})

	cfg := g.ControlFlowGraph(0)
	assert.EqualValues(t, 2, cfg.Size())
	assert.EqualValues(t, [][]uint{{0}, {1, 2}}, cfg.BasicBlockToNodes)
	assert.EqualValues(t, []uint{0, 1, 1}, cfg.NodeToBasicBlock)

	assert.EqualExportedValues(t, control_flow.Node{
		ForwardEdges:  []uint{1},
		BackwardEdges: nil,
	}, cfg.Nodes[0])

	assert.EqualValues(t, control_flow.Node{
		ForwardEdges:  []uint{1},
		BackwardEdges: []uint{1, 0},
	}, cfg.Nodes[1])
}

// func TestLoopToEntry(t *testing.T) {
// 	// input graph:
// 	// 0 -> 1 -> 2
// 	// ^---------+
// 	//
// 	// control flow graph:
// 	// 0-+  (single node, self loop)
// 	// ^-+

// 	instructions := []control_flow.SupportsControlFlow{
// 		&TestInstruction{NextInstructionIndices: []uint{1}},
// 		&TestInstruction{NextInstructionIndices: []uint{2}},
// 		&TestInstruction{NextInstructionIndices: []uint{0}},
// 	}

// 	graph := control_flow.NewControlFlowGraph(instructions)
// 	assert.Len(t, graph.BasicBlocks, 1)

// 	assert.EqualValues(t, control_flow.ControlFlowBasicBlock{
// 		NodeIndices:   []uint{0, 1, 2},
// 		ForwardEdges:  []uint{0},
// 		BackwardEdges: []uint{0},
// 	}, graph.BasicBlocks[0])
// }

// func TestLoopAndJumpOverLoop(t *testing.T) {
// 	// input graph:
// 	// 0 1 2 3
// 	// +-----^ (jump over loop)
// 	//   +-^   (loop)
// 	//   ^-+
// 	//
// 	// control flow graph:
// 	// 0    (two separate components, first component is a singleton, no edges)
// 	// 1-+ (second component, single node, self loop)
// 	// ^-+

// 	instructions := []control_flow.SupportsControlFlow{
// 		&TestInstruction{NextInstructionIndices: []uint{3}}, // 0
// 		&TestInstruction{NextInstructionIndices: []uint{2}}, // 1
// 		&TestInstruction{NextInstructionIndices: []uint{1}}, // 2
// 		&TestInstruction{NextInstructionIndices: []uint{}},  // 3
// 	}

// 	graph := control_flow.NewControlFlowGraph(instructions)
// 	assert.Len(t, graph.BasicBlocks, 2)

// 	assert.EqualValues(t, control_flow.ControlFlowBasicBlock{
// 		NodeIndices:   []uint{0, 3},
// 		ForwardEdges:  []uint{},
// 		BackwardEdges: []uint{},
// 	}, graph.BasicBlocks[0])

// 	assert.EqualValues(t, control_flow.ControlFlowBasicBlock{
// 		NodeIndices:   []uint{1, 2},
// 		ForwardEdges:  []uint{1},
// 		BackwardEdges: []uint{1},
// 	}, graph.BasicBlocks[1])
// }

// func TestIfElse(t *testing.T) {
// 	// input graph:
// 	// 0 1 2 3
// 	// +-^     if
// 	// +---^   else
// 	//   +---^ continue
// 	//     +-^ continue
// 	//
// 	// control flow graph:
// 	//     0
// 	//  v--+--v
// 	//  1     3
// 	//  +--v--+
// 	//     2

// 	instructions := []control_flow.SupportsControlFlow{
// 		&TestInstruction{NextInstructionIndices: []uint{1, 2}}, // 0
// 		&TestInstruction{NextInstructionIndices: []uint{3}},    // 1
// 		&TestInstruction{NextInstructionIndices: []uint{3}},    // 2
// 		&TestInstruction{NextInstructionIndices: []uint{}},     // 3
// 	}

// 	graph := control_flow.NewControlFlowGraph(instructions)
// 	assert.Len(t, graph.BasicBlocks, 4)

// 	// TODO: the following asserts relay on the current implementation and
// 	// assume the order of the basic blocks in the graph. This is not ideal.

// 	assert.EqualValues(t, control_flow.ControlFlowBasicBlock{
// 		NodeIndices:   []uint{0},
// 		ForwardEdges:  []uint{1, 3},
// 		BackwardEdges: []uint{},
// 	}, graph.BasicBlocks[0])

// 	assert.EqualValues(t, control_flow.ControlFlowBasicBlock{
// 		NodeIndices:   []uint{1},
// 		ForwardEdges:  []uint{2},
// 		BackwardEdges: []uint{0},
// 	}, graph.BasicBlocks[1])

// 	assert.EqualValues(t, control_flow.ControlFlowBasicBlock{
// 		NodeIndices:   []uint{3},
// 		ForwardEdges:  []uint{},
// 		BackwardEdges: []uint{1, 3},
// 	}, graph.BasicBlocks[2])

// 	assert.EqualValues(t, control_flow.ControlFlowBasicBlock{
// 		NodeIndices:   []uint{2},
// 		ForwardEdges:  []uint{2},
// 		BackwardEdges: []uint{0},
// 	}, graph.BasicBlocks[3])

// }
