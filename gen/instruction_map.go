package gen

import (
	"strings"

	"alon.kr/x/faststringmap"
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

type InstructionMap struct {
	faststringmap.Map[InstructionDefinition]
	caseSensitive bool
}

func (m *InstructionMap) GetInstructionDefinition(
	name string,
	node parse.InstructionNode,
) (InstructionDefinition, core.ResultList) {
	if !m.caseSensitive {
		name = strings.ToLower(name)
	}

	def, ok := m.Map.LookupString(name)
	if !ok {
		return nil, list.FromSingle(core.Result{{
			Type:     core.ErrorResult,
			Message:  "Undefined instruction",
			Location: &node.Operator,
		}})
	}

	return def, core.ResultList{}
}

func NewInstructionMap(
	definitions []faststringmap.MapEntry[InstructionDefinition],
	caseSensitive bool,
) InstructionManager {
	if !caseSensitive {
		for i := range definitions {
			definitions[i].Key = strings.ToLower(definitions[i].Key)
		}
	}

	return &InstructionMap{
		Map:           faststringmap.NewMap(definitions),
		caseSensitive: caseSensitive,
	}
}
