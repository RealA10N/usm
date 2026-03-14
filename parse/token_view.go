package parse

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/lex"

	"alon.kr/x/view"
)

type TokenView struct{ view.View[lex.Token, uint32] }

func NewTokenView(tkns []lex.Token) TokenView {
	return TokenView{view.NewView[lex.Token, uint32](tkns)}
}

func (v *TokenView) PeekToken(expectedTypes ...lex.TokenType) (tkn lex.Token, result core.Result) {
	tkn, err := v.At(0)

	if err != nil {
		result = NewEofResult(expectedTypes)
		return
	}

	gotExpected := false
	for _, expectedType := range expectedTypes {
		if tkn.Type == expectedType {
			gotExpected = true
			break
		}
	}

	if !gotExpected {
		result = NewUnexpectedTokenResult(expectedTypes, tkn)
		return
	}

	return tkn, nil
}

func (v *TokenView) ConsumeToken(expectedTypes ...lex.TokenType) (tkn lex.Token, result core.Result) {
	tknView, restView := v.Partition(1)
	tkn, err := tknView.At(0)

	if err != nil {
		result = NewEofResult(expectedTypes)
		return
	}

	gotExpected := false
	for _, expectedType := range expectedTypes {
		if tkn.Type == expectedType {
			gotExpected = true
			break
		}
	}

	if !gotExpected {
		result = NewUnexpectedTokenResult(expectedTypes, tkn)
		return
	}

	*v = TokenView{restView}
	return tkn, nil
}

// Consume as many tokens as possible greedily, until we receive an error.
// The error and token count are returned.
func (v *TokenView) ConsumeManyTokens(
	expectedTypes ...lex.TokenType,
) (count int, err core.Result) {
	for ; ; count++ {
		_, err = v.ConsumeToken(expectedTypes...)
		if err != nil {
			return
		}
	}
}

// Consume a token, but ignore any separator or comment tokens that come before it.
func (v *TokenView) ConsumeTokenIgnoreSeparator(
	expectedTypes ...lex.TokenType,
) (lex.Token, core.Result) {
	v.consumeLeadingComments()
	return v.ConsumeToken(expectedTypes...)
}

// consumeLeadingComments consumes any interleaved separators and comment tokens,
// returning the captured comments. This is used to collect whole-line comments
// that appear before a node, so they can be attached to that node.
func (v *TokenView) consumeLeadingComments() []lex.Comment {
	var comments []lex.Comment
	for {
		v.ConsumeManyTokens(lex.SeparatorToken)
		tkn, err := v.ConsumeToken(lex.CommentToken)
		if err != nil {
			return comments
		}
		comments = append(comments, lex.Comment{View: tkn.View})
	}
}

// consumeTrailingComment consumes and returns the next token if it is a
// CommentToken (an inline comment on the same line as the preceding node).
// Returns nil if no comment token is next.
func (v *TokenView) consumeTrailingComment() *lex.Comment {
	tkn, err := v.ConsumeToken(lex.CommentToken)
	if err != nil {
		return nil
	}
	c := lex.Comment{View: tkn.View}
	return &c
}

// Consume as many tokens as possible greedily, until we receive an error.
//
// If the number of tokens consumed is strictly less than the provided number,
// returns the underlying error. Otherwise, returns nil.
func (v *TokenView) ConsumeAtLeastTokens(
	atLeast int, expectedTypes ...lex.TokenType,
) (err core.Result) {
	count, err := v.ConsumeManyTokens(expectedTypes...)
	if count < atLeast {
		return err
	} else {
		return nil
	}
}
