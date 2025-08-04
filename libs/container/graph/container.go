package container_graph

import (
	"fmt"
	"slices"

	"github.com/necroin/golibs/utils"
)

type Graph[T any] struct {
	nodes map[string]*Node[T]
}

func New[T any](nodes ...*Node[T]) *Graph[T] {
	nodeByNames := map[string]*Node[T]{}
	for _, node := range nodes {
		nodeByNames[node.Name()] = node
	}

	return &Graph[T]{
		nodes: nodeByNames,
	}
}

func (container *Graph[T]) Nodes() map[string]*Node[T] {
	return container.nodes
}

func (container *Graph[T]) NodesNames() []string {
	return utils.MapToSlice(
		container.nodes,
		func(key string, node *Node[T]) string { return node.Name() },
	)
}

func (container *Graph[T]) NodesTransitionsNames() map[string][]string {
	return utils.MapToMap(
		container.nodes,
		func(key string, node *Node[T]) (string, []string) { return node.Name(), node.TransitionsNames() },
	)
}

func (container *Graph[T]) GetNode(name string) (*Node[T], error) {
	node, ok := container.nodes[name]
	if !ok {
		return nil, fmt.Errorf("node with name %s not exitst", name)
	}
	return node, nil
}

func (container *Graph[T]) HasNode(name string) bool {
	_, ok := container.nodes[name]
	return ok
}

func (container *Graph[T]) AddNode(name string, value T) (*Node[T], error) {
	node := NewNode(name, value)
	return container.AddNodeItem(node)
}

func (container *Graph[T]) AddNodeItem(node *Node[T]) (*Node[T], error) {
	if container.HasNode(node.Name()) {
		return nil, fmt.Errorf("node with name %s already exists", node.Name())
	}

	container.nodes[node.Name()] = node

	return node, nil
}

func (container *Graph[T]) AddTransition(from, to string, options ...map[string]any) error {
	fromNode := container.nodes[from]
	toNode := container.nodes[to]

	if fromNode == nil {
		return fmt.Errorf("from node do not exist")
	}

	if toNode == nil {
		return fmt.Errorf("to node do not exist")
	}

	fromNode.AddTransition(toNode, options...)

	return nil
}

func (container *Graph[T]) AddTransitionUndirected(n1, n2 string, options ...map[string]any) error {
	err := container.AddTransition(n1, n2, options...)
	if err != nil {
		return err
	}
	return container.AddTransition(n2, n1, options...)
}

func (container *Graph[T]) TopologicalSort() error {
	const (
		unvisited = iota
		visiting
		visited
	)

	states := make(map[*Node[T]]int)
	var result []*Node[T]

	var visit func(*Node[T]) error
	visit = func(node *Node[T]) error {
		switch states[node] {
		case visiting:
			return nil
		case visited:
			return nil
		}

		states[node] = visiting

		for _, neighbor := range node.Transitions() {
			visit(neighbor.node)
		}

		states[node] = visited
		result = append(result, node)
		return nil
	}

	for _, node := range container.nodes {
		if states[node] == unvisited {
			visit(node)
		}
	}

	slices.Reverse(result)
	container.nodes = utils.SliceToMap(result, func(node *Node[T]) (string, *Node[T]) { return node.Name(), node })

	return nil
}
