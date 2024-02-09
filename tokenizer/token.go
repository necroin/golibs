package tokenizer

import (
	"fmt"
	"regexp"
	"strconv"
)

type Token struct {
	name    string
	pattern string
	regex   regexp.Regexp
	value   string
}

func NewToken(name string, pattern string) *Token {
	pattern = "^" + pattern
	regex, err := regexp.Compile(pattern)
	if err != nil {
		panic(fmt.Sprintf("[Token] failed compile %s pattern: %s", pattern, err))
	}
	return &Token{
		name:    name,
		regex:   *regex,
		pattern: pattern,
	}
}

func (token *Token) String() string {
	return fmt.Sprintf("{Name: %s, Pattern: %s}", token.name, token.pattern)
}

func (token *Token) Name() string {
	return token.name
}

func (token *Token) Pattern() string {
	return token.pattern
}

func (token *Token) Value() string {
	return token.value
}

func (token *Token) ValueInt() (int, error) {
	result, err := strconv.Atoi(token.value)
	if err != nil {
		return 0, fmt.Errorf("[Token] [ValueInt] failed convert '%s' value to int: %s", token.value, err)
	}
	return result, nil
}
