package gen

import "alon.kr/x/list"

type BasicBlockInfo[InstT BaseInstruction] struct {
	Instructions list.List[*InstructionInfo[InstT]]

	Labels []*LabelInfo[InstT]
}
