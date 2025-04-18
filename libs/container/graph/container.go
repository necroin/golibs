package container_graph

import (
	"fmt"
	"slices"

	"github.com/necroin/golibs/utils"
)

type Graph[T any] struct {
	nodes       []*Node[T]
	nodeByNames map[string]*Node[T]
}

func New[T any](nodes ...*Node[T]) *Graph[T] {
	nodeByNames := map[string]*Node[T]{}
	for _, node := range nodes {
		nodeByNames[node.name] = node
	}

	return &Graph[T]{
		nodes:       nodes,
		nodeByNames: nodeByNames,
	}
}

func (container *Graph[T]) Nodes() []*Node[T] {
	return container.nodes
}

func (container *Graph[T]) NodesNames() []string {
	return utils.MapSlice(
		container.nodes,
		func(node *Node[T]) string { return node.name },
	)
}

func (container *Graph[T]) NodesTransitionsNames() map[string][]string {
	return utils.SliceToMap(
		container.nodes,
		func(node *Node[T]) (string, []string) { return node.name, node.TransitionsNames() },
	)
}

func (container *Graph[T]) GetNode(name string) (*Node[T], error) {
	node, ok := container.nodeByNames[name]
	if !ok {
		return nil, fmt.Errorf("node with name %s not exitst", name)
	}
	return node, nil
}

func (container *Graph[T]) AddNode(name string, value T) (*Node[T], error) {
	node := NewNode(name, value)
	return container.AddNodeItem(node)
}

func (container *Graph[T]) AddNodeItem(node *Node[T]) (*Node[T], error) {
	container.nodeByNames[node.name] = node
	container.nodes = append(container.nodes, node)
	return node, nil
}

func (container *Graph[T]) AddTransition(from, to string) error {
	fromNode := container.nodeByNames[from]
	toNode := container.nodeByNames[to]

	if fromNode == nil {
		return fmt.Errorf("from node do not exist")
	}

	if toNode == nil {
		return fmt.Errorf("to node do not exist")
	}

	fromNode.AddTransition(toNode)

	return nil
}

func (container *Graph[T]) AddTransitionUndirected(n1, n2 string) error {
	err := container.AddTransition(n1, n2)
	if err != nil {
		return err
	}
	return container.AddTransition(n2, n1)
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

		for _, neighbor := range node.transitions {
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
	container.nodes = result

	return nil
}
