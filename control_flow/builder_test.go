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

func TestBuildEmptyControlFlowGraph(t *testing.T) {
	instructions := []control_flow.SupportsControlFlow{}
	graph := control_flow.NewControlFlowGraph(instructions)
	assert.Empty(t, graph.BasicBlocks)
}
