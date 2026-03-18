package parse

import "alon.kr/x/usm/core"

type ArgumentNode Node

type ArgumentParser struct {
	TypeParser
	RegisterParser  Parser[RegisterNode]
	ImmediateParser *ImmediateParser
	LabelParser     Parser[LabelNode]
	GlobalParser    Parser[GlobalNode]
}

func NewArgumentParser() Parser[ArgumentNode] {
	return ArgumentParser{
		RegisterParser:  NewRegisterParser(),
		ImmediateParser: NewImmediateParser(),
		LabelParser:     NewLabelParser(),
		GlobalParser:    NewGlobalParser(),
	}
}

func (p ArgumentParser) Parse(v *TokenView) (node ArgumentNode, err core.Result) {
	// TODO: make this code neater.

	// Try typed register ($type %name) with backtracking. The type annotation
	// is accepted but discarded — the register already carries its declared type.
	savedV := *v
	if _, typeErr := p.TypeParser.Parse(&savedV); typeErr == nil {
		if regNode, regErr := p.RegisterParser.Parse(&savedV); regErr == nil {
			*v = savedV
			return regNode, nil
		}
	}

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

	var location core.UnmanagedSourceView
	if nextToken, err := v.Front(); err == nil {
		location = nextToken.View
	} else {
		// If there is no tokens left, the location of the error is the end of
		// the source file.
		location = core.NewEofUnmanagedSourceView()
	}

	return nil, core.Result{{
		Type:     core.ErrorResult,
		Message:  "Expected argument",
		Location: &location,
	}}
}
