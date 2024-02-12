package parser

import (
	"fmt"
	"strings"
)

type Rule[T any] struct {
	name    string
	pattern string
	tokens  []string
	handler func(tokens []Token[T]) T
}

func NewRule[T any](name string, pattern string, handler func(tokens []Token[T]) T) *Rule[T] {
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
	for i := 0; i < tokensCount; i++ {
		if rule.tokens[i] != tokens[i].Name() {
			return false
		}
	}
	return true
}

func (rule Rule[T]) String() string {
	return fmt.Sprintf("{%s -> %s}", rule.name, rule.pattern)
}
