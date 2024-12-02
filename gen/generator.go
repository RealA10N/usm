package gen

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

type Generator[InstT BaseInstruction, NodeT parse.Node, InfoT any] interface {
	Generate(
		ctx *GenerationContext,
		node NodeT,
	) (info InfoT, results core.ResultList)
}
