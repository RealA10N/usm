package gen

func ArgumentsToRegisters(
	arguments []ArgumentInfo,
) []*RegisterInfo {
	registers := []*RegisterInfo{}

	for _, arg := range arguments {
		if regArg, ok := arg.(*RegisterArgumentInfo); ok {
			registers = append(registers, regArg.Register)
		}
	}

	return registers
}

func TargetsToRegisters(
	targets []*TargetInfo,
) []*RegisterInfo {
	registers := []*RegisterInfo{}

	for _, target := range targets {
		registers = append(registers, target.Register)
	}

	return registers
}
