package gen_test

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

// MARK: TypesMap

type TypeMap map[string]*gen.TypeInfo

func (t *TypeMap) GetType(name string) *gen.TypeInfo {
	val, ok := (*t)[name]
	if !ok {
		return nil
	}
	return val
}

func (t *TypeMap) NewType(typ *gen.TypeInfo) core.Result {
	if _, exists := (*t)[typ.Name]; exists {
		return &core.GenericResult{
			Type:     core.ErrorResult,
			Message:  "Type already defined",
			Location: &typ.Declaration,
		}
	}

	(*t)[typ.Name] = typ
	return nil
}

func (t *TypeMap) newBuiltinType(name string, size core.UsmUint) core.Result {
	info := &gen.TypeInfo{
		Name:        name,
		Size:        size,
		Declaration: core.UnmanagedSourceView{},
	}
	return t.NewType(info)
}

// MARK: RegisterMap

type RegisterMap map[string]*gen.RegisterInfo

func (r *RegisterMap) GetRegister(name string) *gen.RegisterInfo {
	val, ok := (*r)[name]
	if !ok {
		return nil
	}
	return val
}

func (r *RegisterMap) NewRegister(reg *gen.RegisterInfo) core.Result {
	if _, exists := (*r)[reg.Name]; exists {
		return &core.GenericResult{
			Type:    core.ErrorResult,
			Message: "Register already defined",
		}
	}

	(*r)[reg.Name] = reg
	return nil
}
