package gen

import "alon.kr/x/usm/core"

type FunctionInfo struct {
	Name        string
	Declaration *core.UnmanagedSourceView

	EntryBlock *BasicBlockInfo
	Registers  RegisterManager
	Labels     LabelManager
	Parameters []*RegisterInfo
	Targets    []ReferencedTypeInfo
}

// FunctionInfo can contain information about a function that is defined, or
// only declared that will be linked later.
//
// This method returns true if the function is defined and has a function body,
// and false if it is only declared and has only the signature, parameters and
// targets defined.
func (f *FunctionInfo) IsDefined() bool {
	return f.EntryBlock != nil
}

func (f *FunctionInfo) CollectBasicBlocks() []*BasicBlockInfo {
	blocks := make([]*BasicBlockInfo, 0)
	for block := f.EntryBlock; block != nil; block = block.NextBlock {
		blocks = append(blocks, block)
	}

	return blocks
}

func (f *FunctionInfo) CollectInstructions() []*InstructionInfo {
	instructions := make([]*InstructionInfo, 0)
	for block := f.EntryBlock; block != nil; block = block.NextBlock {
		instructions = append(instructions, block.Instructions...)
	}

	return instructions
}

// Returns the number of instructions in the function.
func (f *FunctionInfo) Size() int {
	return len(f.CollectInstructions())
}

func (i *FunctionInfo) String() string {
	s := "func "
	for _, target := range i.Targets {
		s += target.String() + " "
	}

	s += i.Name + " "

	for _, param := range i.Parameters {
		// TODO: create a separate ParameterInfo type and just call String()
		// on it.
		s += param.Type.String() + " " + param.String() + " "
	}

	if i.IsDefined() {
		s += "{\n"

		for block := i.EntryBlock; block != nil; block = block.NextBlock {
			s += block.String()
		}

		s += "}"
	}

	return s
}
