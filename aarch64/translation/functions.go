package aarch64translation

import (
	"bytes"
	"encoding/binary"

	aarch64instructions "alon.kr/x/aarch64codegen/instructions"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

func FunctionToBinaryData(function *gen.FunctionInfo) ([]byte, core.ResultList) {
	results := core.ResultList{}
	instructions := function.CollectInstructions()
	aarch64Instructions := make(
		[]aarch64instructions.Instruction,
		0,
		len(instructions),
	)

	for _, instruction := range instructions {
		aarchInstruction, ok := instruction.Instruction.(aarch64instructions.Instruction)
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
		return nil, results
	}

	data := bytes.Buffer{}
	for _, instruction := range aarch64Instructions {
		binary.Write(&data, binary.LittleEndian, instruction.Binary())
	}

	return data.Bytes(), core.ResultList{}
}
