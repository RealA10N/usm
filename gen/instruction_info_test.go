package gen_test

import (
	"testing"

	"alon.kr/x/usm/gen"
	"github.com/stretchr/testify/assert"
)

func newTestInstruction() *gen.InstructionInfo {
	return gen.NewEmptyInstructionInfo(nil)
}

func newTestRegister(name string) *gen.RegisterInfo {
	return gen.NewRegisterInfo(name, gen.ReferencedTypeInfo{})
}

// TestSubstituteArgument_RegisterToRegister verifies that swapping one register
// argument for another updates the References lists on both RegisterInfo structs.
func TestSubstituteArgument_RegisterToRegister(t *testing.T) {
	oldReg := newTestRegister("%old")
	newReg := newTestRegister("%new")
	instruction := newTestInstruction()

	oldArg := gen.NewRegisterArgumentInfo(oldReg)
	instruction.AppendArgument(oldArg)

	newArg := gen.NewRegisterArgumentInfo(newReg)
	instruction.SubstituteArgument(0, newArg)

	assert.NotContains(t, oldReg.References, instruction, "old register should no longer list the instruction as a reference")
	assert.Contains(t, newReg.References, instruction, "new register should list the instruction as a reference")
	assert.Equal(t, newArg, instruction.Arguments[0])
}

// TestSubstituteArgument_RegisterToImmediate verifies that replacing a register
// argument with an immediate removes the instruction from the register's References.
func TestSubstituteArgument_RegisterToImmediate(t *testing.T) {
	reg := newTestRegister("%x")
	instruction := newTestInstruction()

	regArg := gen.NewRegisterArgumentInfo(reg)
	instruction.AppendArgument(regArg)

	imm := &gen.ImmediateInfo{}
	instruction.SubstituteArgument(0, imm)

	assert.NotContains(t, reg.References, instruction, "register should no longer list the instruction as a reference")
	assert.Equal(t, imm, instruction.Arguments[0])
}

// TestSubstituteArgument_ImmediateToRegister verifies that replacing an
// immediate argument with a register adds the instruction to the register's References.
func TestSubstituteArgument_ImmediateToRegister(t *testing.T) {
	reg := newTestRegister("%x")
	instruction := newTestInstruction()

	instruction.AppendArgument(&gen.ImmediateInfo{})

	regArg := gen.NewRegisterArgumentInfo(reg)
	instruction.SubstituteArgument(0, regArg)

	assert.Contains(t, reg.References, instruction, "register should list the instruction as a reference after substitution")
	assert.Equal(t, regArg, instruction.Arguments[0])
}

// TestSubstituteArgument_MultipleArgs verifies that SubstituteArgument only
// modifies the References of the register at the specified index, not others.
func TestSubstituteArgument_MultipleArgs(t *testing.T) {
	reg0 := newTestRegister("%a")
	reg1 := newTestRegister("%b")
	newReg := newTestRegister("%c")
	instruction := newTestInstruction()

	arg0 := gen.NewRegisterArgumentInfo(reg0)
	arg1 := gen.NewRegisterArgumentInfo(reg1)
	instruction.AppendArgument(arg0, arg1)

	newArg := gen.NewRegisterArgumentInfo(newReg)
	instruction.SubstituteArgument(1, newArg)

	assert.Contains(t, reg0.References, instruction, "reg0 at index 0 should be unaffected")
	assert.NotContains(t, reg1.References, instruction, "reg1 at index 1 should be removed")
	assert.Contains(t, newReg.References, instruction, "newReg should be added")
}

// TestSwitchRegister verifies that SwitchRegister updates both References lists
// and that the argument's Register field is updated in place (so the
// instruction's Arguments slice sees the new register without re-slicing).
func TestSwitchRegister(t *testing.T) {
	oldReg := newTestRegister("%old")
	newReg := newTestRegister("%new")
	instruction := newTestInstruction()

	arg := gen.NewRegisterArgumentInfo(oldReg)
	instruction.AppendArgument(arg)

	arg.SwitchRegister(instruction, newReg)

	assert.NotContains(t, oldReg.References, instruction)
	assert.Contains(t, newReg.References, instruction)
	// The argument object in the slice is the same pointer; its Register field
	// has been updated in place.
	assert.Equal(t, newReg, instruction.Arguments[0].(*gen.RegisterArgumentInfo).Register)
}
