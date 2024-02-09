package tokenizer

import (
	"bytes"
	"fmt"
	"sort"
)

type Tokenizer struct {
	tokens       []*Token
	values       []*Token
	ignoreSpaces bool
}

func NewTokenizer(tokens ...*Token) *Tokenizer {
	return &Tokenizer{
		tokens:       tokens,
		values:       []*Token{},
		ignoreSpaces: true,
	}
}

func (tokenizer *Tokenizer) Tokens() []*Token {
	return tokenizer.values
}

func (tokenizer *Tokenizer) sortFindTokens() {
	sort.Slice(tokenizer.tokens, func(i, j int) bool {
		return len(tokenizer.tokens[i].pattern) > len(tokenizer.tokens[j].pattern)
	})
}

func (tokenizer *Tokenizer) Find(text []byte) (*Token, error) {
	for _, token := range tokenizer.tokens {
		findedValue := token.regex.Find(text)
		if len(findedValue) != 0 {
			valuedToken := &Token{
				name:    token.name,
				pattern: token.pattern,
				value:   string(findedValue),
			}
			return valuedToken, nil
		}
	}
	return nil, fmt.Errorf("[Tokenizer] no tokens matched")
}

func (tokenizer *Tokenizer) Parse(text []byte) error {
	tokenizer.sortFindTokens()
	tokenizer.values = []*Token{}

	for len(text) != 0 {
		if tokenizer.ignoreSpaces {
			text = bytes.Trim(text, " ")
		}

		token, err := tokenizer.Find(text)
		if err != nil {
			return fmt.Errorf("[Tokenizer] [Parse] failed parse text: %s", err)
		}

		tokenizer.values = append(tokenizer.values, token)
		text = text[len(token.Value()):]
	}

	return nil
}

func (tokenizer *Tokenizer) SetIgnoreSpaces(value bool) {
	tokenizer.ignoreSpaces = value
}
