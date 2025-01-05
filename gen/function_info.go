package gen

import (
	"alon.kr/x/set"
	"alon.kr/x/stack"
)

type FunctionInfo struct {
	Name       string
	EntryBlock *BasicBlockInfo
	Registers  RegisterManager
	Labels     LabelManager
	Parameters []*RegisterInfo
	Targets    []ReferencedTypeInfo
}

// Performs a DFS traversal of the function's basic blocks, and returns a set
// of the collected blocks in the traversal.
func (f *FunctionInfo) CollectBasicBlocks() []*BasicBlockInfo {
	blocks := make([]*BasicBlockInfo, 0)
	visited := set.New[*BasicBlockInfo]()
	toVisit := stack.New[*BasicBlockInfo]()

	blocks = append(blocks, f.EntryBlock)
	toVisit.Push(f.EntryBlock)
	visited.Add(f.EntryBlock)

	for len(toVisit) > 0 {
		block := toVisit.Top()
		toVisit.Pop()

		for _, nextBlock := range block.ForwardEdges {
			if !visited.Contains(nextBlock) {
				blocks = append(blocks, nextBlock)
				toVisit.Push(nextBlock)
				visited.Add(nextBlock)
			}
		}
	}

	return blocks
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

	s += "{\n"

	for _, block := range i.CollectBasicBlocks() {
		s += block.String()
	}

	s += "}"
	return s
}
