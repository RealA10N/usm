package parse

type ArgumentNode Node

type ArgumentParser struct {
	RegisterParser  RegisterParser
	ImmediateParser *ImmediateParser
	LabelParser     Parser[LabelNode]
	GlobalParser    GlobalParser
}

func NewArgumentParser() ArgumentParser {
	return ArgumentParser{
		ImmediateParser: NewImmediateParser(),
		LabelParser:     NewLabelParser(),
	}
}

func (p ArgumentParser) Parse(v *TokenView) (node ArgumentNode, err ParsingError) {
	// TODO: make this code neater.

	if node, err := p.RegisterParser.Parse(v); err == nil {
		return node, nil
	}

	if node, err := p.ImmediateParser.Parse(v); err == nil {
		return node, nil
	}

	if node, err := p.LabelParser.Parse(v); err == nil {
		return node, nil
	}

	// TODO: make global part of immediate?
	if node, err := p.GlobalParser.Parse(v); err == nil {
		return node, nil
	}

	return nil, GenericUnexpectedError{Expected: "argument"}
}
