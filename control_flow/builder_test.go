package control_flow_test

import (
	"testing"

	"alon.kr/x/usm/control_flow"
	"github.com/stretchr/testify/assert"
)

type TestInstruction struct {
	NextInstructionIndices []uint
}

func (i *TestInstruction) PossibleNextInstructionIndices() []uint {
	return i.NextInstructionIndices
}

func TestBuildEmpty(t *testing.T) {
	instructions := []control_flow.SupportsControlFlow{}
	graph := control_flow.NewControlFlowGraph(instructions)
	assert.Empty(t, graph.BasicBlocks)
}

func TestSingleBlock(t *testing.T) {
	instructions := []control_flow.SupportsControlFlow{
		&TestInstruction{NextInstructionIndices: []uint{1}},
		&TestInstruction{NextInstructionIndices: []uint{2}},
		&TestInstruction{NextInstructionIndices: []uint{}},
	}

	graph := control_flow.NewControlFlowGraph(instructions)
	assert.Len(t, graph.BasicBlocks, 1)
	assert.Equal(t, []uint{0, 1, 2}, graph.BasicBlocks[0].InstructionIndices)
	assert.Empty(t, graph.BasicBlocks[0].ForwardEdges)
}

func TestSimpleLoop(t *testing.T) {
	instructions := []control_flow.SupportsControlFlow{
		&TestInstruction{NextInstructionIndices: []uint{1}},
		&TestInstruction{NextInstructionIndices: []uint{2}},
		&TestInstruction{NextInstructionIndices: []uint{1}},
	}

	graph := control_flow.NewControlFlowGraph(instructions)
	assert.Len(t, graph.BasicBlocks, 2)

	assert.EqualValues(t, control_flow.ControlFlowBasicBlock{
		InstructionIndices: []uint{0},
		ForwardEdges:       []uint{1},
	}, graph.BasicBlocks[0])

	assert.EqualValues(t, control_flow.ControlFlowBasicBlock{
		InstructionIndices: []uint{1, 2},
		ForwardEdges:       []uint{1},
	}, graph.BasicBlocks[1])
}

func TestLoopToEntry(t *testing.T) {
	instructions := []control_flow.SupportsControlFlow{
		&TestInstruction{NextInstructionIndices: []uint{1}},
		&TestInstruction{NextInstructionIndices: []uint{2}},
		&TestInstruction{NextInstructionIndices: []uint{0}},
	}

	graph := control_flow.NewControlFlowGraph(instructions)
	assert.Len(t, graph.BasicBlocks, 1)

	assert.EqualValues(t, control_flow.ControlFlowBasicBlock{
		InstructionIndices: []uint{0, 1, 2},
		ForwardEdges:       []uint{0},
	}, graph.BasicBlocks[0])
}

func TestLoopAndJumpOverLoop(t *testing.T) {
	// 0 1 2 3
	// +-----^ jump over loop
	//   +-^   loop
	//   ^-+

	instructions := []control_flow.SupportsControlFlow{
		&TestInstruction{NextInstructionIndices: []uint{3}}, // 0
		&TestInstruction{NextInstructionIndices: []uint{2}}, // 1
		&TestInstruction{NextInstructionIndices: []uint{1}}, // 2
		&TestInstruction{NextInstructionIndices: []uint{}},  // 3
	}

	graph := control_flow.NewControlFlowGraph(instructions)
	assert.Len(t, graph.BasicBlocks, 2)

	assert.EqualValues(t, control_flow.ControlFlowBasicBlock{
		InstructionIndices: []uint{0, 3},
		ForwardEdges:       []uint{},
	}, graph.BasicBlocks[0])

	assert.EqualValues(t, control_flow.ControlFlowBasicBlock{
		InstructionIndices: []uint{1, 2},
		ForwardEdges:       []uint{1},
	}, graph.BasicBlocks[1])
}

func TestIfElse(t *testing.T) {
	// 0 1 2 3
	// +-^     if
	// +---^   else
	//   +---^ continue
	//     +-^ continue

	instructions := []control_flow.SupportsControlFlow{
		&TestInstruction{NextInstructionIndices: []uint{1, 2}}, // 0
		&TestInstruction{NextInstructionIndices: []uint{3}},    // 1
		&TestInstruction{NextInstructionIndices: []uint{3}},    // 2
		&TestInstruction{NextInstructionIndices: []uint{}},     // 3
	}

	graph := control_flow.NewControlFlowGraph(instructions)
	assert.Len(t, graph.BasicBlocks, 4)

	// TODO: the following asserts relay on the current implementation and
	// assume the order of the basic blocks in the graph. This is not ideal.

	assert.EqualValues(t, control_flow.ControlFlowBasicBlock{
		InstructionIndices: []uint{0},
		ForwardEdges:       []uint{1, 3},
	}, graph.BasicBlocks[0])

	assert.EqualValues(t, control_flow.ControlFlowBasicBlock{
		InstructionIndices: []uint{1},
		ForwardEdges:       []uint{2},
	}, graph.BasicBlocks[1])

	assert.EqualValues(t, control_flow.ControlFlowBasicBlock{
		InstructionIndices: []uint{3},
		ForwardEdges:       []uint{},
	}, graph.BasicBlocks[2])

	assert.EqualValues(t, control_flow.ControlFlowBasicBlock{
		InstructionIndices: []uint{2},
		ForwardEdges:       []uint{2},
	}, graph.BasicBlocks[3])

}
