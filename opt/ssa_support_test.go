package opt_test

import (
	"math/big"
	"testing"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/opt"
	"github.com/stretchr/testify/assert"
)

// buildTestInstruction creates a minimal InstructionInfo with one target register
// (%a) and one argument register (%b), for use in mixin method tests.
func buildTestInstruction() (*gen.InstructionInfo, *gen.RegisterArgumentInfo, *gen.RegisterArgumentInfo) {
	intType := gen.NewNamedTypeInfo("$32", big.NewInt(32), nil)
	typeRef := gen.ReferencedTypeInfo{Base: intType}
	v := core.UnmanagedSourceView{}

	regA := &gen.RegisterInfo{Name: "%a", Type: typeRef, Declaration: v}
	regB := &gen.RegisterInfo{Name: "%b", Type: typeRef, Declaration: v}

	argA := gen.NewRegisterArgumentInfo(regA)
	argB := gen.NewRegisterArgumentInfo(regB)

	info := gen.NewEmptyInstructionInfo(&v)
	info.AppendTarget(argA)
	info.AppendArgument(argB)

	return info, argA, argB
}

func TestUsesArgumentsInstruction(t *testing.T) {
	info, _, argB := buildTestInstruction()
	args := opt.UsesArgumentsInstruction{}.UsesArguments(info)
	assert.Len(t, args, 1)
	assert.Equal(t, argB, args[0])
}

func TestUsesNothingInstruction(t *testing.T) {
	info, _, _ := buildTestInstruction()
	uses := opt.UsesNothingInstruction{}
	assert.Empty(t, uses.Uses(info))
	assert.Empty(t, uses.UsesArguments(info))
}

func TestDefinesTargetsInstructionDefinitionArguments(t *testing.T) {
	info, argA, _ := buildTestInstruction()
	defArgs := opt.DefinesTargetsInstruction{}.DefinitionArguments(info)
	assert.Len(t, defArgs, 1)
	assert.Equal(t, argA, defArgs[0])
}

func TestDefinesNothingInstruction(t *testing.T) {
	info, _, _ := buildTestInstruction()
	defines := opt.DefinesNothingInstruction{}
	assert.Empty(t, defines.Defines(info))
	assert.Empty(t, defines.DefinitionArguments(info))
}
