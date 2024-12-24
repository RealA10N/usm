package managers

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	usm64core "alon.kr/x/usm/usm64/core"
)

type LabelMap map[string]*gen.LabelInfo[usm64core.Instruction]

func (m *LabelMap) GetLabel(name string) *gen.LabelInfo[usm64core.Instruction] {
	val, ok := (*m)[name]
	if !ok {
		return nil
	}
	return val
}

func (m *LabelMap) NewLabel(label *gen.LabelInfo[usm64core.Instruction]) core.Result {
	(*m)[label.Name] = label
	return nil
}

func NewLabelManager() gen.LabelManager[usm64core.Instruction] {
	return &LabelMap{}
}
