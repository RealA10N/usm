package parse

type ArgumentNode Node

type ArgumentParser struct {
	RegisterParser  RegisterParser
	ImmediateParser ImmediateParser
	LabelParser     LabelParser
	GlobalParser    GlobalParser
}

func (ArgumentParser) String() string {
	return "argument"
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

	if node, err := p.GlobalParser.Parse(v); err == nil {
		return node, nil
	}

	return nil, GenericUnexpectedError{Expected: p.String()}
}
