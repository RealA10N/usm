package gen

type FunctionInfo[InstT BaseInstruction] struct {
	Instructions []*InstructionInfo[InstT]
	Parameters   []*RegisterInfo
	// TODO: add targets
}
