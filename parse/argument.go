package parse

import (
	"usm/source"
)

type ArgumentNode struct {
	Type     TypeNode
	Register RegisterNode
}

func (n ArgumentNode) View() source.UnmanagedSourceView {
	return n.Type.View().Merge(n.Register.View())
}

type ArgumentParser struct{}

func (ArgumentParser) Parse(v *TokenView) (node ArgumentNode, err ParsingError) {
	typ, err := TypeParser{}.Parse(v)
	if err != nil {
		return
	}

	reg, err := RegisterParser{}.Parse(v)
	if err != nil {
		return
	}

	node = ArgumentNode{
		Type:     typ,
		Register: reg,
	}

	return
}
