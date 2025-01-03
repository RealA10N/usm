package gen

import (
	"strconv"

	"alon.kr/x/usm/core"
)

type TypeDescriptorType uint8

const (
	PointerTypeDescriptor TypeDescriptorType = iota
	RepeatTypeDescriptor
)

func (t TypeDescriptorType) String() string {
	switch t {
	case PointerTypeDescriptor:
		return "*"
	case RepeatTypeDescriptor:
		return "^"
	default:
		panic("unreachable")
	}
}

type TypeDescriptorInfo struct {
	Type   TypeDescriptorType
	Amount core.UsmUint
}

func (i TypeDescriptorInfo) String() string {
	amount := i.Amount
	strconv.Itoa(int(i.Amount)) // TODO: remove type conversion. Use bigint instead?
	return i.Type.String() + strconv.Itoa(int(amount))
}

// A referenced type is a combination of a basic type with (possibly zero)
// type decorators that wrap it.
// For example, if `$32“ is a basic named type, then `$32 *`, which is a
// pointer to that type is a referenced type with the `$32` named type as it's
// base type, and the pointer as a decorator.
type ReferencedTypeInfo struct {
	// A pointer to the base, named type that this type reference refers to.
	Base        *NamedTypeInfo
	Descriptors []TypeDescriptorInfo
}

func (t ReferencedTypeInfo) String() string {
	s := t.Base.String()
	for _, descriptor := range t.Descriptors {
		s += " " + descriptor.String()
	}

	return s
}

func (info ReferencedTypeInfo) Equal(other ReferencedTypeInfo) bool {
	if info.Base != other.Base {
		return false
	}

	if len(info.Descriptors) != len(other.Descriptors) {
		return false
	}

	for i := range info.Descriptors {
		if info.Descriptors[i] != other.Descriptors[i] {
			return false
		}
	}

	return true
}
