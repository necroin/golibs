package parser

type Token[T any] interface {
	Name() string
	Value() T
	String() string
}

type ParserToken[T any] struct {
	name  string
	value T
}

func NewParserToken[T any](name string, value T) *ParserToken[T] {
	return &ParserToken[T]{
		name:  name,
		value: value,
	}
}

func (token ParserToken[T]) Name() string {
	return token.name
}

func (token ParserToken[T]) Value() T {
	return token.value
}

func (token ParserToken[T]) String() string {
	return token.name
}
