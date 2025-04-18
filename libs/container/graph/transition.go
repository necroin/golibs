package container_graph

type Transition[T any] struct {
	node    *Node[T]
	options map[string]any
}

func (transition *Transition[T]) Node() *Node[T] {
	return transition.node
}

func (transition *Transition[T]) Options() map[string]any {
	return transition.options
}
