package gen

import "alon.kr/x/usm/core"

// This represents partial register information, possibly without an associated
// type (yet). This is used internally before the compiler has finally determined
// the type of the register, if the type is implicit.
type registerPartialInfo struct {
	Name string

	// Possibly nil, if type is implicitly defined.
	Type *ReferencedTypeInfo

	Declaration core.UnmanagedSourceView
}
