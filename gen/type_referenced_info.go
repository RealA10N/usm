package gen

import (
	"math/big"

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
	Amount *big.Int
}

func (i TypeDescriptorInfo) String() string {
	return i.Type.String() + i.Amount.String()
}

func (i TypeDescriptorInfo) Equal(other TypeDescriptorInfo) bool {
	return i.Type == other.Type && i.Amount.Cmp(other.Amount) == 0
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

	Declaration *core.UnmanagedSourceView
}

func (t ReferencedTypeInfo) String() string {
	s := t.Base.String()
	for _, descriptor := range t.Descriptors {
		s += " " + descriptor.String()
	}

	return s
}

func (t ReferencedTypeInfo) IsPure() bool {
	return len(t.Descriptors) == 0
}

// PointerTo returns a new type that is one pointer level deeper than t.
// Because *N packs N pointer levels into a single descriptor, adding one
// level to a type ending with *N yields *(N+1); otherwise a new *1 is appended.
// The returned type has no Declaration since it is derived, not from source.
func (t ReferencedTypeInfo) PointerTo() ReferencedTypeInfo {
	n := len(t.Descriptors)
	newDescriptors := make([]TypeDescriptorInfo, n)
	copy(newDescriptors, t.Descriptors)

	if n > 0 && t.Descriptors[n-1].Type == PointerTypeDescriptor {
		last := newDescriptors[n-1]
		newDescriptors[n-1] = TypeDescriptorInfo{
			Type:   PointerTypeDescriptor,
			Amount: new(big.Int).Add(last.Amount, big.NewInt(1)),
		}
	} else {
		newDescriptors = append(newDescriptors, TypeDescriptorInfo{
			Type:   PointerTypeDescriptor,
			Amount: big.NewInt(1),
		})
	}

	return ReferencedTypeInfo{Base: t.Base, Descriptors: newDescriptors}
}

// Deref returns a new type with one pointer level removed and true, or the
// zero value and false if t does not end with a pointer descriptor.
// The returned type has no Declaration since it is derived, not from source.
func (t ReferencedTypeInfo) Deref() (ReferencedTypeInfo, bool) {
	n := len(t.Descriptors)
	if n == 0 || t.Descriptors[n-1].Type != PointerTypeDescriptor {
		return ReferencedTypeInfo{}, false
	}

	last := t.Descriptors[n-1]
	newDescriptors := make([]TypeDescriptorInfo, n)
	copy(newDescriptors, t.Descriptors)

	if last.Amount.Cmp(big.NewInt(1)) == 0 {
		newDescriptors = newDescriptors[:n-1]
	} else {
		newDescriptors[n-1] = TypeDescriptorInfo{
			Type:   PointerTypeDescriptor,
			Amount: new(big.Int).Sub(last.Amount, big.NewInt(1)),
		}
	}

	return ReferencedTypeInfo{Base: t.Base, Descriptors: newDescriptors}, true
}

func (info ReferencedTypeInfo) Equal(other ReferencedTypeInfo) bool {
	if info.Base != other.Base {
		return false
	}

	if len(info.Descriptors) != len(other.Descriptors) {
		return false
	}

	for i := range info.Descriptors {
		if !info.Descriptors[i].Equal(other.Descriptors[i]) {
			return false
		}
	}

	return true
}
