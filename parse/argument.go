package parse

type ArgumentNode Node

type ArgumentParser struct {
	RegisterParser  Parser[RegisterNode]
	ImmediateParser *ImmediateParser
	LabelParser     Parser[LabelNode]
	GlobalParser    Parser[GlobalNode]
}

func NewArgumentParser() ArgumentParser {
	return ArgumentParser{
		RegisterParser:  NewRegisterParser(),
		ImmediateParser: NewImmediateParser(),
		LabelParser:     NewLabelParser(),
		GlobalParser:    NewGlobalParser(),
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
