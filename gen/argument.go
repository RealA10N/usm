package gen

import "alon.kr/x/usm/parse"

// MARK: Info

type ArgumentInfo struct {
	// A pointer to the TypeInfo instance that corresponds to the type of the
	// register.
	Type *TypeInfo
}

// MARK: Generator

type ArgumentGenerator[InstT BaseInstruction] Generator[InstT, parse.ArgumentNode, *ArgumentInfo]
