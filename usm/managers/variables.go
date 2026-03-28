package usmmanagers

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type VariableMap map[string]*gen.VariableInfo

func (m *VariableMap) GetVariable(name string) *gen.VariableInfo {
	val, ok := (*m)[name]
	if !ok {
		return nil
	}
	return val
}

func (m *VariableMap) NewVariable(variable *gen.VariableInfo) core.ResultList {
	(*m)[variable.Name] = variable
	return core.ResultList{}
}

func (m *VariableMap) GetAllVariables() []*gen.VariableInfo {
	variables := make([]*gen.VariableInfo, 0, len(*m))
	for _, v := range *m {
		variables = append(variables, v)
	}
	return variables
}

func NewVariableManager(*gen.FileGenerationContext) gen.VariableManager {
	return &VariableMap{}
}
