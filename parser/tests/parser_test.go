package tests

import (
	"fmt"
	"testing"

	"github.com/necroin/golibs/parser"
)

type TestToken struct {
	name  string
	value string
}

func NewTestToken(name string, value string) *TestToken {
	return &TestToken{
		name:  name,
		value: value,
	}
}

func (token TestToken) Name() string {
	return token.name
}

func (token TestToken) Value() string {
	return token.value
}

func (token TestToken) String() string {
	return token.name
}

func TestParser(t *testing.T) {
	parserInstance := parser.NewParser[int]()

	parserInstance.AddRule(parser.NewRule[int]("EXPR", "NUMBER", func(tokens []parser.Token[int]) int {
		return tokens[0].Value()
	}))
	parserInstance.AddRule(parser.NewRule[int]("EXPR", "EXPR OPERATOR EXPR", func(tokens []parser.Token[int]) int {
		return tokens[0].Value() + tokens[2].Value()
	}))

	result, err := parserInstance.Parse(
		parser.ParseOptions{
			LogFunc: func(format string, args ...any) { fmt.Printf(format+"\n", args...) },
		},
		parser.NewParserToken[int]("NUMBER", 5),
		parser.NewParserToken[int]("OPERATOR", 1),
		parser.NewParserToken[int]("NUMBER", 5),
	)
	if err != nil {
		t.Fatal(err)
	}

	if result.Value() != 10 {
		t.Fatal("Wrong result")
	}
}

func TestParserCustom(t *testing.T) {
	parserInstance := parser.NewParser[string]()

	parserInstance.AddRule(parser.NewRule[string]("EXPR", "TEXT", func(tokens []parser.Token[string]) string {
		return tokens[0].Value()
	}))
	parserInstance.AddRule(parser.NewRule[string]("EXPR", "EXPR OPERATOR EXPR", func(tokens []parser.Token[string]) string {
		return tokens[0].Value() + tokens[2].Value()
	}))

	result, err := parserInstance.Parse(
		parser.ParseOptions{
			LogFunc: func(format string, args ...any) { fmt.Printf(format+"\n", args...) },
		},
		NewTestToken("TEXT", "Hello"),
		NewTestToken("OPERATOR", "+"),
		NewTestToken("TEXT", " World"),
	)
	if err != nil {
		t.Fatal(err)
	}

	if result.Value() != "Hello World" {
		t.Fatal("Wrong result")
	}
}
