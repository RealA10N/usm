package gen_test

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

// MARK: TypesMap

type TypeMap map[string]*gen.NamedTypeInfo

func (m *TypeMap) GetType(name string) *gen.NamedTypeInfo {
	val, ok := (*m)[name]
	if !ok {
		return nil
	}
	return val
}

func (m *TypeMap) NewType(typ *gen.NamedTypeInfo) core.Result {
	(*m)[typ.Name] = typ
	return nil
}

func (m *TypeMap) newBuiltinType(name string, size core.UsmUint) core.Result {
	info := &gen.NamedTypeInfo{
		Name:        name,
		Size:        size,
		Declaration: core.UnmanagedSourceView{},
	}
	return m.NewType(info)
}

// MARK: RegisterMap

type RegisterMap map[string]*gen.RegisterInfo

func (m *RegisterMap) GetRegister(name string) *gen.RegisterInfo {
	val, ok := (*m)[name]
	if !ok {
		return nil
	}
	return val
}

func (m *RegisterMap) NewRegister(reg *gen.RegisterInfo) core.Result {
	(*m)[reg.Name] = reg
	return nil
}

// MARK: LabelMap

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

// MARK: Context

var testInstructionSet = gen.InstructionManager[Instruction](
	&InstructionMap{
		"ADD": &AddInstructionDefinition{},
	},
)

var testManagerCreators = gen.ManagerCreators{
	LabelManagerCreator: func() gen.LabelManager {
		return gen.LabelManager(&LabelMap{})
	},
	RegisterManagerCreator: func() gen.RegisterManager {
		return gen.RegisterManager(&RegisterMap{})
	},
	TypeManagerCreator: func() gen.TypeManager {
		return gen.TypeManager(&TypeMap{})
	},
}

var testGenerationContext = gen.GenerationContext[Instruction]{
	ManagerCreators: testManagerCreators,
	Instructions:    testInstructionSet,
	PointerSize:     314, // An arbitrary, unique value.
}
