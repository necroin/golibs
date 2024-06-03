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
	ignoreTabs   bool
}

func NewTokenizer(tokens ...*Token) *Tokenizer {
	tokens = tokens[:]

	sort.Slice(tokens, func(i, j int) bool {
		return len(tokens[i].pattern) > len(tokens[j].pattern)
	})

	return &Tokenizer{
		tokens:       tokens,
		ignoreSpaces: true,
		ignoreTabs:   true,
	}
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
	return nil, fmt.Errorf("[Tokenizer] no tokens matched: %s", text)
}

func (tokenizer *Tokenizer) Parse(text []byte) ([]*Token, error) {
	tokens := []*Token{}

	for len(text) != 0 {
		trimCutset := ""

		if tokenizer.ignoreSpaces {
			trimCutset = trimCutset + " "
		}

		if tokenizer.ignoreTabs {
			trimCutset = trimCutset + "\t"
		}

		text = bytes.Trim(text, trimCutset)

		token, err := tokenizer.Find(text)
		if err != nil {
			return nil, fmt.Errorf("[Tokenizer] [Parse] failed parse text: %s", err)
		}

		tokens = append(tokens, token)
		text = text[len(token.Value()):]
	}

	return tokens, nil
}

func (tokenizer *Tokenizer) SetIgnoreSpaces(value bool) {
	tokenizer.ignoreSpaces = value
}

func (tokenizer *Tokenizer) SetIgnoreTabs(value bool) {
	tokenizer.ignoreTabs = value
}
