package aarch64translation

import (
	"bytes"
	"encoding/binary"

	aarch64codegen "alon.kr/x/usm/aarch64/codegen"
	aarch64core "alon.kr/x/usm/aarch64/core"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

func FunctionCodegen(
	fileCtx *aarch64codegen.FileCodegenContext,
	function *gen.FunctionInfo,
	buffer *bytes.Buffer,
) core.ResultList {
	results := core.ResultList{}
	instructions := function.CollectInstructions()
	aarch64Instructions := make(
		[]aarch64core.Instruction,
		0,
		len(instructions),
	)

	for _, instruction := range instructions {
		aarchInstruction, ok := instruction.Instruction.(aarch64core.Instruction)
		if ok {
			aarch64Instructions = append(aarch64Instructions, aarchInstruction)
		} else {
			results.Append(core.Result{
				{
					Type:     core.InternalErrorResult,
					Message:  "Instruction is not an AArch64 instruction",
					Location: instruction.Declaration,
				},
			})
		}
	}

	if !results.IsEmpty() {
		return results
	}

	funcCtx := aarch64codegen.NewFunctionCodegenContext(fileCtx, function)
	for _, instruction := range aarch64Instructions {
		binaryInst := instruction.Generate(funcCtx)
		binary.Write(buffer, binary.LittleEndian, binaryInst.Binary())
	}

	return core.ResultList{}
}
