package managers

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type LabelMap map[string]gen.LabelInfo

func (m *LabelMap) GetLabel(name string) *gen.LabelInfo {
	val, ok := (*m)[name]
	if !ok {
		return nil
	}
	return &val
}

func (m *LabelMap) NewLabel(label gen.LabelInfo) core.Result {
	(*m)[label.Name] = label
	return nil
}

func NewLabelManager() gen.LabelManager {
	return &LabelMap{}
}
