package parse_test

import (
	"testing"
	"usm/parse"
	"usm/source"

	"github.com/stretchr/testify/assert"
)

func TestInstructionStringer(t *testing.T) {
	inst := "%div %mod = divmod %x %y"
	v, ctx := source.NewSourceView(inst).Detach()

	node := parse.InstructionNode{
		Operator: v.Subview(12, 18),
		Arguments: []parse.ArgumentNode{
			parse.ArgumentNode{v.Subview(19, 21)},
			parse.ArgumentNode{v.Subview(22, 24)},
		},
		Targets: []parse.RegisterNode{
			parse.RegisterNode{v.Subview(0, 4)},
			parse.RegisterNode{v.Subview(5, 9)},
		},
	}

	assert.Equal(t, inst, node.String(ctx))
}
