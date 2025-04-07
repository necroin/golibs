package container_graph

import (
	"fmt"
	"slices"

	"github.com/necroin/golibs/utils"
)

type Node[T any] struct {
	name        string
	value       T
	transitions []*Node[T]
	options     map[string]any
}

func NewNode[T any](name string, value T) *Node[T] {
	return &Node[T]{
		name:        name,
		value:       value,
		transitions: []*Node[T]{},
		options:     map[string]any{},
	}
}

func (node *Node[T]) Name() string {
	return node.name
}

func (node *Node[T]) Value() T {
	return node.value
}

func (node *Node[T]) Transitions() []*Node[T] {
	return node.transitions
}

func (node *Node[T]) Options() map[string]any {
	return node.options
}

func (node *Node[T]) SetOption(name string, value any) *Node[T] {
	node.options[name] = value
	return node
}

func (node *Node[T]) TransitionsNames() []string {
	return utils.MapSlice(node.transitions, func(node *Node[T]) string { return node.name })
}

func (node *Node[T]) AddTransition(toNode *Node[T]) {
	if node == toNode {
		return
	}

	if slices.Contains(node.transitions, toNode) {
		return
	}

	node.transitions = append(node.transitions, toNode)
}

func (node *Node[T]) String() string {
	return fmt.Sprintf("{name: %s, value: %v}", node.name, node.value)
}
