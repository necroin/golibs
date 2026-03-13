package parser

import (
	"fmt"
	"strings"
)

type RuleHandler[T any] func(tokens []Token[T]) T

type Rule[T any] struct {
	name    string
	pattern string
	tokens  []string
	handler RuleHandler[T]
}

func NewRule[T any](name string, pattern string, handler RuleHandler[T]) *Rule[T] {
	if len(name) == 0 {
		panic("[Parser] [Rule] empty name")
	}
	if len(pattern) == 0 {
		panic("[Parser] [Rule] empty pattern")
	}
	return &Rule[T]{
		name:    name,
		pattern: pattern,
		tokens:  strings.Split(pattern, " "),
		handler: handler,
	}
}

func (rule Rule[T]) CompareTokens(tokens []Token[T]) bool {
	tokensCount := len(rule.tokens)
	for i := range tokensCount {
		if rule.tokens[i] != tokens[i].Name() {
			return false
		}
	}
	return true
}

func (rule Rule[T]) String() string {
	return fmt.Sprintf("{%s -> %s}", rule.pattern, rule.name)
}
