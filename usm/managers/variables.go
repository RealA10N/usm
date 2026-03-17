package usmmanagers

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

// VariableList is the default implementation of gen.VariableManager.
// It preserves declaration order via a slice and uses a map for O(1) lookup.
type VariableList struct {
	byName []*gen.VariableInfo
	index  map[string]int
}

func (m *VariableList) GetVariable(name string) *gen.VariableInfo {
	i, ok := m.index[name]
	if !ok {
		return nil
	}
	return m.byName[i]
}

func (m *VariableList) NewVariable(variable *gen.VariableInfo) core.ResultList {
	m.index[variable.Name] = len(m.byName)
	m.byName = append(m.byName, variable)
	return core.ResultList{}
}

func (m *VariableList) GetAllVariables() []*gen.VariableInfo {
	result := make([]*gen.VariableInfo, len(m.byName))
	copy(result, m.byName)
	return result
}

func NewVariableManager(*gen.FileGenerationContext) gen.VariableManager {
	return &VariableList{
		index: make(map[string]int),
	}
}
