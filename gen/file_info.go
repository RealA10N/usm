package gen

type FileInfo[InstT BaseInstruction] struct {
	Functions []*FunctionInfo[InstT]
}
