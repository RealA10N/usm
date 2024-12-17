package control_flow

type SupportsControlFlow interface {
	// Returns a list of all instruction indices in the function that execution
	// could arrive to after the execution of this instruction.
	PossibleNextInstructionIndices() []uint
}

type ControlFlowBasicBlock struct {
	InstructionIndices []uint
	ForwardEdges       []uint
	BackwardEdges      []uint
}

type ControlFlowGraph[InstT SupportsControlFlow] struct {
	Instructions []InstT
	BasicBlocks  []ControlFlowBasicBlock
}
