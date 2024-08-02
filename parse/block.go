package parse

import (
	"usm/lex"
	"usm/source"
)

// TODO: add label support

type BlockNode struct {
	source.UnmanagedSourceView
	Instructions []InstructionNode
}

func (n BlockNode) View() source.UnmanagedSourceView {
	return n.UnmanagedSourceView
}

func (n BlockNode) String(ctx source.SourceContext) (s string) {
	s = "{\n"
	for _, inst := range n.Instructions {
		s += "\t" + inst.String(ctx) + "\n"
	}
	s += "}\n"
	return s
}

type BlockParser struct{}

func (p BlockParser) Parse(v *TokenView) (node BlockNode, err ParsingError) {
	start, err := v.ConsumeToken(lex.LcrToken)
	if err != nil {
		return
	}

	for {
		inst, err := InstructionParser{}.Parse(v)
		if err != nil {
			break
		}
		node.Instructions = append(node.Instructions, inst)
	}

	end, err := v.ConsumeToken(lex.RcrToken)
	if err != nil {
		return
	}

	node.UnmanagedSourceView = start.View.Merge(end.View)
	return node, nil
}
