package gen

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

// MARK: Info
type RegisterInfo struct {
	// The name of the register, as it appears in the source code.
	Name string

	// The type of the register.
	Type ReferencedTypeInfo

	// The first location in the source code in which the register is declared
	// or assigned a value.
	Declaration core.UnmanagedSourceView
}

func (i RegisterInfo) toPartialRegisterInfo() partialRegisterInfo {
	return partialRegisterInfo{
		Name:        i.Name,
		Type:        &i.Type,
		Declaration: i.Declaration,
	}
}

// This represents partial register information, possibly without an associated
// type (yet). This is used internally before the compiler has finally determined
// the type of the register, if the type is implicit.
type partialRegisterInfo struct {
	Name string

	// Possibly nil, if type is implicitly defined.
	Type *ReferencedTypeInfo

	Declaration core.UnmanagedSourceView
}

// Converts the partial register information type into a full register information
// structure, with the a guaranteed register type.
//
// Returns an error if the provided actual type does not match the explicit
// partial type.
func (i partialRegisterInfo) toRegisterInfo(
	actualType ReferencedTypeInfo,
) (RegisterInfo, core.ResultList) {
	if i.Type != nil && !i.Type.Equals(actualType) {
		return RegisterInfo{}, list.FromSingle(core.Result{
			{
				Type:     core.InternalErrorResult,
				Message:  "Explicit type does not match implicit type",
				Location: &i.Declaration,
			},
		})
	}

	info := RegisterInfo{
		Name:        i.Name,
		Type:        actualType,
		Declaration: i.Declaration,
	}

	return info, core.ResultList{}
}

// MARK: Manager

type RegisterManager interface {
	GetRegister(name string) *RegisterInfo
	NewRegister(reg *RegisterInfo) core.Result
}

// MARK: Generator

// Used to convert parse.RegisterNode instances to *existing* register instances.
// Returns an error on generation if the provided register node references a
// register that does not exist.
type RegisterGenerator[InstT BaseInstruction] struct{}

func NewRegisterGenerator[InstT BaseInstruction]() Generator[InstT, parse.RegisterNode, *RegisterInfo] {
	return Generator[InstT, parse.RegisterNode, *RegisterInfo](
		&RegisterGenerator[InstT]{},
	)
}

func UndefinedRegisterResult(node parse.RegisterNode) core.ResultList {
	v := node.View()
	return list.FromSingle(core.Result{
		{
			Type:     core.ErrorResult,
			Message:  "Undefined register",
			Location: &v,
		},
	})
}

func getRegisterNameFromRegisterNode[InstT BaseInstruction](
	ctx *GenerationContext[InstT],
	node parse.RegisterNode,
) string {
	return string(node.Raw(ctx.SourceContext))
}

func (g *RegisterGenerator[InstT]) Generate(
	ctx *GenerationContext[InstT],
	node parse.RegisterNode,
) (*RegisterInfo, core.ResultList) {
	name := getRegisterNameFromRegisterNode(ctx, node)
	register := ctx.Registers.GetRegister(name)

	if register == nil {
		return nil, UndefinedRegisterResult(node)
	}

	return register, core.ResultList{}
}
