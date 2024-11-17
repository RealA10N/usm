package core

import (
	"strconv"
	"unsafe"
)

// Base type for unsigned integers.
// This defines, for example, the maximum length of a source file, and the
// maximum type size, and other bounds for internal data structures.
type UsmUint = uint32

const UsmUintBitSize = 8 * unsafe.Sizeof(UsmUint(0))

func ParseUint(s string) (UsmUint, error) {
	n, err := strconv.ParseUint(s, 10, int(UsmUintBitSize))
	return UsmUint(n), err
}
