package parse_test

import (
	"testing"

	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/parse"
	"alon.kr/x/usm/source"

	"github.com/stretchr/testify/assert"
)

func TestSignatureParserOnlyIdentifier(t *testing.T) {
	v, ctx := source.NewSourceView("@fibonacci").Detach()
	glb := lex.Token{Type: lex.GlobalToken, View: v}
	tknView := parse.NewTokenView([]lex.Token{glb})

	expectedSig := parse.FunctionSignatureNode{
		UnmanagedSourceView: v.Subview(0, 10),
		Identifier:          v,
	}

	sig, err := parse.NewFunctionSignatureParser().Parse(&tknView)
	assert.Nil(t, err)
	assert.Equal(t, expectedSig, sig)

	assert.Equal(t, v, sig.View())
	strCtx := parse.StringContext{SourceContext: ctx}
	assert.Equal(t, "@fibonacci", sig.String(&strCtx))
}

func TestSignatureParserVoidFunction(t *testing.T) {
	v, ctx := source.NewSourceView("... @printNumber  $i32 %x ...").Detach()
	glb := lex.Token{Type: lex.GlobalToken, View: v.Subview(4, 16)}
	typ := lex.Token{Type: lex.TypeToken, View: v.Subview(18, 22)}
	reg := lex.Token{Type: lex.RegisterToken, View: v.Subview(23, 25)}
	tknView := parse.NewTokenView([]lex.Token{glb, typ, reg})

	expectedSig := parse.FunctionSignatureNode{
		UnmanagedSourceView: v.Subview(4, 25),
		Identifier:          glb.View,
		Parameters: []parse.ParameterNode{
			{
				Type:     parse.TypeNode{Identifier: typ.View},
				Register: parse.RegisterNode{reg.View},
			},
		},
	}

	sig, err := parse.NewFunctionSignatureParser().Parse(&tknView)
	assert.Nil(t, err)
	assert.Equal(t, expectedSig, sig)
	assert.Equal(t, v.Subview(4, 25), sig.View())
	strCtx := parse.StringContext{SourceContext: ctx}
	assert.Equal(t, "@printNumber $i32 %x", sig.String(&strCtx))
}

func TestSignatureParserSingleReturn(t *testing.T) {
	v, ctx := source.NewSourceView("$i32 @add $i32 %x $i32 %y").Detach()
	ret := lex.Token{Type: lex.TypeToken, View: v.Subview(0, 4)}
	id := lex.Token{Type: lex.GlobalToken, View: v.Subview(5, 9)}
	param1Typ := lex.Token{Type: lex.TypeToken, View: v.Subview(10, 14)}
	param1Reg := lex.Token{Type: lex.RegisterToken, View: v.Subview(15, 17)}
	param2Typ := lex.Token{Type: lex.TypeToken, View: v.Subview(18, 22)}
	param2Reg := lex.Token{Type: lex.RegisterToken, View: v.Subview(23, 25)}
	tknView := parse.NewTokenView([]lex.Token{
		ret, id, param1Typ, param1Reg, param2Typ, param2Reg,
	})

	expectedSig := parse.FunctionSignatureNode{
		UnmanagedSourceView: v,
		Identifier:          id.View,
		Parameters: []parse.ParameterNode{
			{
				Type:     parse.TypeNode{Identifier: param1Typ.View},
				Register: parse.RegisterNode{param1Reg.View},
			},
			{
				Type:     parse.TypeNode{Identifier: param2Typ.View},
				Register: parse.RegisterNode{param2Reg.View},
			},
		},
		Returns: []parse.TypeNode{
			{Identifier: ret.View},
		},
	}

	sig, err := parse.NewFunctionSignatureParser().Parse(&tknView)
	assert.Nil(t, err)
	assert.Equal(t, expectedSig, sig)
	assert.Equal(t, v, sig.View())
	strCtx := parse.StringContext{SourceContext: ctx}
	assert.Equal(t, "$i32 @add $i32 %x $i32 %y", sig.String(&strCtx))
}

func TestSignatureParserMutltiReturn(t *testing.T) {
	v, ctx := source.NewSourceView("$i32  $i32  @divmod  $i32 %n  $i32 %m").Detach()
	ret1 := lex.Token{Type: lex.TypeToken, View: v.Subview(0, 4)}
	ret2 := lex.Token{Type: lex.TypeToken, View: v.Subview(6, 10)}
	id := lex.Token{Type: lex.GlobalToken, View: v.Subview(12, 19)}
	param1Typ := lex.Token{Type: lex.TypeToken, View: v.Subview(21, 25)}
	param1Reg := lex.Token{Type: lex.RegisterToken, View: v.Subview(26, 28)}
	param2Typ := lex.Token{Type: lex.TypeToken, View: v.Subview(30, 34)}
	param2Reg := lex.Token{Type: lex.RegisterToken, View: v.Subview(35, 37)}
	tknView := parse.NewTokenView([]lex.Token{
		ret1, ret2, id, param1Typ, param1Reg, param2Typ, param2Reg,
	})

	expectedSig := parse.FunctionSignatureNode{
		UnmanagedSourceView: v,
		Identifier:          id.View,
		Parameters: []parse.ParameterNode{
			{
				Type:     parse.TypeNode{Identifier: param1Typ.View},
				Register: parse.RegisterNode{param1Reg.View},
			},
			{
				Type:     parse.TypeNode{Identifier: param2Typ.View},
				Register: parse.RegisterNode{param2Reg.View},
			},
		},
		Returns: []parse.TypeNode{
			{Identifier: ret1.View},
			{Identifier: ret2.View},
		},
	}

	sig, err := parse.NewFunctionSignatureParser().Parse(&tknView)
	assert.Nil(t, err)
	assert.Equal(t, expectedSig, sig)
	assert.Equal(t, v, sig.View())
	strCtx := parse.StringContext{SourceContext: ctx}
	assert.Equal(t, "$i32 $i32 @divmod $i32 %n $i32 %m", sig.String(&strCtx))
}

func TestSignatureParserIdentifierNotGlobal(t *testing.T) {
	v, _ := source.NewSourceView("funcid").Detach()
	opr := lex.Token{Type: lex.OperatorToken, View: v}
	tknView := parse.NewTokenView([]lex.Token{opr})

	expectedErr := parse.UnexpectedTokenError{
		Expected: []lex.TokenType{lex.GlobalToken},
		Actual:   opr,
	}

	_, err := parse.NewFunctionSignatureParser().Parse(&tknView)
	assert.NotNil(t, err)

	details, ok := err.(parse.UnexpectedTokenError)
	assert.True(t, ok)
	assert.Equal(t, expectedErr, details)
}
