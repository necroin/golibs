package tests

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/necroin/golibs/parser"
	"github.com/necroin/golibs/tokenizer"
)

type ExpressionData struct {
	IntValue    int64
	StringValue string
}

func (ed *ExpressionData) String() string {
	return fmt.Sprintf("{IntValue: %d, StringValue: %s}", ed.IntValue, ed.StringValue)
}

func Calculate(left int64, right int64, operator string) int64 {
	switch operator {
	case "+":
		return left + right
	case "-":
		return left - right
	case "*":
		return left * right
	case "/":
		return left / right
	}
	return 0
}

func TestExpression(t *testing.T) {
	refData := map[string]int64{
		"ref1":  10,
		"ref_2": 20,
	}

	expression := "10 * $(ref1) + $(ref_2) - (20 / 10)"

	tokenizer := tokenizer.NewTokenizer(
		tokenizer.NewToken("REF_START", `\$`),
		tokenizer.NewToken("OPEN_BRACKET", `\(`),
		tokenizer.NewToken("CLOSE_BRACKET", `\)`),
		tokenizer.NewToken("REF", `[a-zA-Z_][a-zA-Z0-9_]*`),
		tokenizer.NewToken("OPERATOR", `[\+\-\*\/]`),
		tokenizer.NewToken("NUMBER", `[0-9]+`),
	)

	tokens, err := tokenizer.Parse([]byte(expression))
	if err != nil {
		t.Fatal(err)
	}

	expressionParser := parser.NewParser[ExpressionData]()

	expressionParser.AddRule(parser.NewRule[ExpressionData]("EXPR", "EXPR OPERATOR EXPR", func(tokens []parser.Token[ExpressionData]) ExpressionData {
		return ExpressionData{
			IntValue: Calculate(tokens[0].Value().IntValue, tokens[2].Value().IntValue, tokens[1].Value().StringValue),
		}
	}))

	expressionParser.AddRule(parser.NewRule[ExpressionData]("EXPR", "OPEN_BRACKET EXPR CLOSE_BRACKET", func(tokens []parser.Token[ExpressionData]) ExpressionData {
		return ExpressionData{
			IntValue: tokens[1].Value().IntValue,
		}
	}))

	expressionParser.AddRule(parser.NewRule[ExpressionData]("EXPR", "NUMBER", func(tokens []parser.Token[ExpressionData]) ExpressionData {
		result, _ := strconv.ParseInt(tokens[0].Value().StringValue, 10, 64)
		return ExpressionData{
			IntValue: result,
		}
	}))

	expressionParser.AddRule(parser.NewRule[ExpressionData]("EXPR", "REF_START OPEN_BRACKET REF CLOSE_BRACKET", func(tokens []parser.Token[ExpressionData]) ExpressionData {
		return ExpressionData{
			IntValue: refData[tokens[2].Value().StringValue],
		}
	}))

	parserTokens := []parser.Token[ExpressionData]{}
	for _, token := range tokens {
		parserTokens = append(parserTokens, parser.NewParserToken(token.Name(), ExpressionData{StringValue: token.Value()}))
	}

	_, err = expressionParser.Parse(
		parser.ParseOptions{
			LogFunc: func(format string, args ...any) { fmt.Printf(format+"\n", args...) },
		},
		parserTokens...,
	)
	if err != nil {
		t.Fatal(err)
	}
}
