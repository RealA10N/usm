package gen

import "alon.kr/x/usm/graph"

type FunctionInfo[InstT BaseInstruction] struct {

	// A linear view of the instructions in the program.
	graph.Graph
	Instructions []*InstructionInfo[InstT]

	// A Control Flow Graph representation of the program.
	graph.ControlFlowGraph
	BasicBlocks []*BasicBlockInfo[InstT]

	Parameters []*RegisterInfo
	// TODO: add targets
}
