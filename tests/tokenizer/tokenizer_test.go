package tokenizer_tests

import (
	"fmt"
	"testing"

	"github.com/necroin/golibs/libs/tokenizer"
)

func TestTokenizer(t *testing.T) {
	expression := "10 * $(ref1) + $(ref_2) - 20 / 10"

	tokenizer := tokenizer.NewTokenizer(
		tokenizer.NewToken("REF_START", `\$`),
		tokenizer.NewToken("OPEN_BRACKET", `\(`),
		tokenizer.NewToken("CLOSE_BRACKET", `\)`),
		tokenizer.NewToken("REF", `[a-zA-Z_][a-zA-Z0-9_]*`),
		tokenizer.NewToken("OPERATOR_PLUS", `\+`),
		tokenizer.NewToken("OPERATOR_MINUS", `\-`),
		tokenizer.NewToken("OPERATOR_MUL", `\*`),
		tokenizer.NewToken("OPERATOR_DIV", `\/`),
		tokenizer.NewToken("NUMBER", `[0-9]+`),
	)

	parsedTokens, err := tokenizer.Parse([]byte(expression))
	if err != nil {
		t.Fatal(err)
	}

	expectedTokens := []string{
		"NUMBER",
		"OPERATOR_MUL",
		"REF_START",
		"OPEN_BRACKET",
		"REF",
		"CLOSE_BRACKET",
		"OPERATOR_PLUS",
		"REF_START",
		"OPEN_BRACKET",
		"REF",
		"CLOSE_BRACKET",
		"OPERATOR_MINUS",
		"NUMBER",
		"OPERATOR_DIV",
		"NUMBER",
	}

	if len(expectedTokens) != len(parsedTokens) {
		t.Fatalf("Wrong tokens count: %v != %v", parsedTokens, expectedTokens)
	}

	fmt.Println(parsedTokens)

	for i, expectedToken := range expectedTokens {
		parsedToken := parsedTokens[i]
		if expectedToken != parsedToken.Name() {
			t.Fatalf("Failed compare tokens: %v != %v", parsedTokens, expectedTokens)
		}
	}
}
