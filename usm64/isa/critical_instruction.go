package usm64isa

// A critical instruction is an instruction that cannot be removed by the
// dead code elimination process. This means that it has side effects, or
// is a function call, or is a branch, etc.
type CriticalInstruction struct{}

func (i *CriticalInstruction) IsCritical() bool {
	return true
}
