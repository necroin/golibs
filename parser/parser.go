package parser

import (
	"fmt"
	"sort"
)

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

func (parser *Parser[T]) Parse(tokens ...Token[T]) error {
	fmt.Println(parser.rules)
	parser.sortRules()
	fmt.Println(parser.rules)

	offset := 0
	for len(tokens) > 1 {
		fmt.Printf("Iteration tokens: %s\n", tokens)
		fmt.Printf("Iteration offset: %d\n", offset)
		if offset == len(tokens) {
			return fmt.Errorf("failed parse tokens: %s", tokens)
		}
		matched := false
		for _, rule := range parser.rules {
			fmt.Printf("Rule: %s\n", rule)

			ruleTokensCount := len(rule.tokens)
			tokensCount := len(tokens)
			if ruleTokensCount+offset > tokensCount {
				continue
			}
			matchTokens := tokens[offset : ruleTokensCount+offset]
			fmt.Printf("Match tokens: %s\n", matchTokens)
			if rule.CompareTokens(matchTokens) {
				fmt.Printf("Reduce by rule: %s\n", rule)
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
	fmt.Printf("Final tokens: %s\n", tokens)
	fmt.Printf("Result: %v\n", tokens[0].Value())

	return nil
}
