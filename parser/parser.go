package parser

import (
	"fmt"
	"sort"
)

type ParseOptions struct {
	LogFunc func(format string, args ...any)
}

type Parser[T any] struct {
	rules []*Rule[T]
}

func NewParser[T any](rules ...*Rule[T]) *Parser[T] {
	return &Parser[T]{
		rules: rules,
	}
}

func (parser *Parser[T]) AddRule(rule *Rule[T]) {
	parser.rules = append(parser.rules, rule)
}

func (parser *Parser[T]) sortRules() {
	sort.Slice(parser.rules, func(i, j int) bool {
		leftTokensCount := len(parser.rules[i].tokens)
		rightTokensCount := len(parser.rules[j].tokens)
		return leftTokensCount > rightTokensCount
	})
}

func (parser *Parser[T]) Parse(options ParseOptions, tokens ...Token[T]) (Token[T], error) {
	parser.sortRules()
	if options.LogFunc != nil {
		options.LogFunc("[Parser] parse rules: %s", parser.rules)
	}

	offset := 0
	matched := false

	for len(tokens) > 1 || matched {
		if options.LogFunc != nil {
			options.LogFunc("[Parser] iteration tokens: %s", tokens)
			options.LogFunc("[Parser] iteration offset: %d", offset)
		}
		if offset == len(tokens) {
			return nil, fmt.Errorf("[Parser] failed parse tokens: %s", tokens)
		}
		matched = false
		for _, rule := range parser.rules {
			if options.LogFunc != nil {
				options.LogFunc("[Parser] rule: %s", rule)
			}

			ruleTokensCount := len(rule.tokens)
			tokensCount := len(tokens)
			if offset+ruleTokensCount > tokensCount {
				if options.LogFunc != nil {
					options.LogFunc("[Parser] decline rule: offset+ruleTokensCount (%d) > tokensCount (%d)", offset+ruleTokensCount, tokensCount)
				}
				continue
			}
			matchTokens := tokens[offset : ruleTokensCount+offset]
			if options.LogFunc != nil {
				options.LogFunc("[Parser] match tokens: %s", matchTokens)
			}
			if rule.CompareTokens(matchTokens) {
				if options.LogFunc != nil {
					options.LogFunc("[Parser] reduce by rule: %s", rule)
				}
				newTokens := append(tokens[:offset], NewParserToken[T](rule.name, rule.handler(matchTokens)))
				newTokens = append(newTokens, tokens[ruleTokensCount+offset:]...)
				tokens = newTokens
				offset = 0
				matched = true
				break
			}
		}
		if !matched {
			offset += 1
		}
	}
	if options.LogFunc != nil {
		options.LogFunc("[Parser] final tokens: %s", tokens)
		fmt.Printf("[Parser] result token value: %v", tokens[0].Value())
	}

	return tokens[0], nil
}
