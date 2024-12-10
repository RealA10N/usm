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
	n, err := strconv.ParseUint(s, 0, int(UsmUintBitSize))
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
