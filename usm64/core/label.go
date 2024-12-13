package usm64core

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type Label struct {
	Name             string
	InstructionIndex uint64
	declaration      core.UnmanagedSourceView
}

func NewLabel(arg gen.LabelArgumentInfo) (Label, core.ResultList) {
	return Label{
		Name:             arg.Label.Name,
		InstructionIndex: arg.Label.InstructionIndex,
		declaration:      arg.Declaration(),
	}, core.ResultList{}
}

func (l Label) String(*EmulationContext) string {
	return l.Name
}

func (l Label) Declaration() core.UnmanagedSourceView {
	return l.declaration
}
