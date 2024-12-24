package gen

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

// Used to convert parse.RegisterNode instances to *existing* register instances.
// Returns an error on generation if the provided register node references a
// register that does not exist.
type RegisterGenerator[InstT BaseInstruction] struct{}

func NewRegisterGenerator[InstT BaseInstruction]() FunctionContextGenerator[InstT, parse.RegisterNode, *RegisterInfo] {
	return FunctionContextGenerator[InstT, parse.RegisterNode, *RegisterInfo](
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

func (g *RegisterGenerator[InstT]) Generate(
	ctx *FunctionGenerationContext[InstT],
	node parse.RegisterNode,
) (*RegisterInfo, core.ResultList) {
	name := nodeToSourceString(ctx.FileGenerationContext, node)
	register := ctx.Registers.GetRegister(name)

	if register == nil {
		return nil, UndefinedRegisterResult(node)
	}

	return register, core.ResultList{}
}
