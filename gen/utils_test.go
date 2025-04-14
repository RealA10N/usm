package gen_test

import (
	"fmt"
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

func (m *RegisterMap) Size() uint {
	return uint(len(*m))
}

func (m *RegisterMap) GetAllRegisters() []*gen.RegisterInfo {
	registers := make([]*gen.RegisterInfo, 0, len(*m))
	for _, reg := range *m {
		registers = append(registers, reg)
	}
	return registers
}

// MARK: LabelMap

type LabelMap map[string]*gen.LabelInfo

func (m *LabelMap) GetLabel(name string) *gen.LabelInfo {
	val, ok := (*m)[name]
	if !ok {
		return nil
	}
	return val
}

func (m *LabelMap) NewLabel(label *gen.LabelInfo) core.ResultList {
	(*m)[label.Name] = label
	return core.ResultList{}
}

func (m *LabelMap) GenerateLabel() *gen.LabelInfo {
	return &gen.LabelInfo{
		Name: ".L" + fmt.Sprint(len(*m)),
	}
}

// MARK: Context

var testInstructionSet = gen.InstructionManager(
	&InstructionMap{
		"ADD": &AddInstructionDefinition{},
		"JMP": &JumpInstructionDefinition{},
		"JZ":  &JumpZeroInstructionDefinition{},
		"RET": &RetInstructionDefinition{},
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

var testGenerationContext = gen.GenerationContext{
	ManagerCreators: testManagerCreators,
	Instructions:    testInstructionSet,
	PointerSize:     big.NewInt(314), // An arbitrary, unique value.
}
