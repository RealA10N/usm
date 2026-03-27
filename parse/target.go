package parse

// TargetNode represents a typed register reference in instruction target
// position.  It is structurally identical to RegisterNode (an optional type
// annotation followed by a register name), so it is defined as a type alias.
type TargetNode = RegisterNode

// NewTargetParser returns a parser for target nodes.  Because TargetNode is
// now an alias for RegisterNode, and RegisterParser already handles the
// optional "$type %reg" syntax, no separate target parser is needed.
func NewTargetParser() Parser[TargetNode] {
	return NewRegisterParser()
}
