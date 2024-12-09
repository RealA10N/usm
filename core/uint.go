package core

import (
	"math/bits"
	"strconv"
	"unsafe"
)

// Base type for unsigned integers.
// This defines, for example, the maximum length of a source file, and the
// maximum type size, and other bounds for internal data structures.
type UsmUint = uint32

const UsmUintBitSize = 8 * unsafe.Sizeof(UsmUint(0))

func ParseUint(s string) (UsmUint, error) {
	// TODO: assert that strconv.ParseUint does not have edge cases and is dead
	// simple, so it won't break our grammar and allow weird behavior which is
	// not defined in the USM spec.
	n, err := strconv.ParseUint(s, 10, int(UsmUintBitSize))
	return UsmUint(n), err
}

// MARK: Arithmetics
// Implementation is based on math.bits.
// Only a subset of required operations is implemented (add more as needed).

func Add(x, y UsmUint) (res UsmUint, ok bool) {
	res, carry := bits.Add32(x, y, 0)
	return res, carry == 0
}

func Mul(x, y UsmUint) (res UsmUint, ok bool) {
	high, low := bits.Mul32(x, y)
	return low, high == 0
}
