package gen_test

import (
	"math/big"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

// MARK: TypeMap

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

func (m *TypeMap) newBuiltinType(name string, size *big.Int) core.Result {
	info := gen.NewNamedTypeInfo(name, size, nil)
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

func (m *RegisterMap) NewRegister(reg *gen.RegisterInfo) core.ResultList {
	(*m)[reg.Name] = reg
	return core.ResultList{}
}

func (m *RegisterMap) DeleteRegister(register *gen.RegisterInfo) core.ResultList {
	delete(*m, register.Name)
	return core.ResultList{}
}

func (m *RegisterMap) Size() int {
	return len(*m)
}

func (m *RegisterMap) GetAllRegisters() []*gen.RegisterInfo {
	registers := make([]*gen.RegisterInfo, 0, len(*m))
	for _, reg := range *m {
		registers = append(registers, reg)
	}
	return registers
}

// MARK: Context

var testInstructionSet = gen.InstructionManager(
	&InstructionMap{
		"add": NewAdd(),
		"j":   NewJump(),
		"jz":  NewJumpZero(),
		"ret": NewRet(),
	},
)

var testManagerCreators = gen.ManagerCreators{
	LabelManagerCreator: gen.NewLabelMap,
	RegisterManagerCreator: func(*gen.FileGenerationContext) gen.RegisterManager {
		return gen.RegisterManager(&RegisterMap{})
	},
	TypeManagerCreator: func(*gen.GenerationContext) gen.TypeManager {
		manager := gen.TypeManager(&TypeMap{})
		typ := gen.NewNamedTypeInfo("$32", big.NewInt(32), nil)
		manager.NewType(typ)
		return manager
	},
	GlobalManagerCreator: gen.NewGlobalMap,
}

var testGenerationContext = gen.GenerationContext{
	ManagerCreators: testManagerCreators,
	Instructions:    testInstructionSet,
	PointerSize:     big.NewInt(314), // An arbitrary, unique value.
}
