package container_graph

import (
	"fmt"
	"slices"

	"github.com/necroin/golibs/utils"
)

type Node[T any] struct {
	name        string
	value       T
	transitions []*Transition[T]
	options     map[string]any
}

func NewNode[T any](name string, value T) *Node[T] {
	return &Node[T]{
		name:        name,
		value:       value,
		transitions: []*Transition[T]{},
		options:     map[string]any{},
	}
}

func (node *Node[T]) Name() string {
	return node.name
}

func (node *Node[T]) Value() T {
	return node.value
}

func (node *Node[T]) Transitions() []*Transition[T] {
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
	return utils.MapSlice(node.transitions, func(transition *Transition[T]) string { return transition.node.name })
}

func (node *Node[T]) AddTransition(toNode *Node[T], options ...map[string]any) {
	if node == toNode {
		return
	}

	if slices.ContainsFunc(node.transitions, func(transition *Transition[T]) bool { return transition.node == toNode }) {
		return
	}

	transitionOptions := map[string]any{}

	for _, option := range options {
		for key, value := range option {
			transitionOptions[key] = value
		}
	}

	node.transitions = append(node.transitions, &Transition[T]{
		node:    toNode,
		options: transitionOptions,
	})
}

func (node *Node[T]) String() string {
	return fmt.Sprintf("{name: %s, value: %v}", node.name, node.value)
}
