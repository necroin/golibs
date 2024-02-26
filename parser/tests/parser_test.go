package tests

import (
	"fmt"
	"testing"

	"github.com/necroin/golibs/parser"
	"github.com/necroin/golibs/tokenizer"
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

	if err := parserInstance.Parse(
		parser.NewParserToken[int]("NUMBER", 5),
		parser.NewParserToken[int]("OPERATOR", 1),
		parser.NewParserToken[int]("NUMBER", 5),
	); err != nil {
		t.Fatal(err)
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

	if err := parserInstance.Parse(
		NewTestToken("TEXT", "Hello"),
		NewTestToken("OPERATOR", "+"),
		NewTestToken("TEXT", " World"),
	); err != nil {
		t.Fatal(err)
	}
}

func TestGraphQl(t *testing.T) {
	expression := `
		{
			HID
			GUID
		}
	`

	tokenizer := tokenizer.NewTokenizer(
		tokenizer.NewToken("QUERY", "query"),
		tokenizer.NewToken("OPEN_BRACE", `\{`),
		tokenizer.NewToken("CLOSE_BRACE", `\}`),
		tokenizer.NewToken("OPEN_BRACKET", `\(`),
		tokenizer.NewToken("CLOSE_BRACKET", `\)`),
		tokenizer.NewToken("WORD", `[a-zA-Z_][a-zA-Z0-9_]*`),
		tokenizer.NewToken("NEWLINE", "\n"),
	)

	if err := tokenizer.Parse([]byte(expression)); err != nil {
		t.Fatal(err)
	}

	tokens := tokenizer.Tokens()
	tokens = tokens[1:]
	for _, token := range tokens {
		fmt.Printf("%s: %s\n", token.Name(), token.Value())
	}

	parserInstance := parser.NewParser[string]()

	parserInstance.AddRule(parser.NewRule[string]("REQUEST_BODY_FIELD", "OPEN_BRACE FIELDS CLOSE_BRACE", func(tokens []parser.Token[string]) string {
		return ""
	}))

	parserInstance.AddRule(parser.NewRule[string]("FIELDS", "FIELDS FIELD", func(tokens []parser.Token[string]) string {
		return ""
	}))

	parserInstance.AddRule(parser.NewRule[string]("FIELDS", "FIELD", func(tokens []parser.Token[string]) string {
		return ""
	}))

	parserInstance.AddRule(parser.NewRule[string]("FIELD", "WORD", func(tokens []parser.Token[string]) string {
		return ""
	}))

	parserTokens := []parser.Token[string]{}
	for _, token := range tokens {
		if token.Name() != "NEWLINE" {
			parserTokens = append(parserTokens, token)
		}
	}

	if err := parserInstance.Parse(parserTokens...); err != nil {
		t.Fatal(err)
	}
}
