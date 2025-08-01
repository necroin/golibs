package container_graph

import (
	"fmt"
	"strings"
)

func (container *Graph[T]) VisualizeMermaid() string {
	var builder strings.Builder

	builder.WriteString("```mermaid\n")
	builder.WriteString("graph LR\n")

	for _, node := range container.nodes {
		for _, transition := range node.transitions {
			builder.WriteString(fmt.Sprintf("\t%s --> %s\n", node.name, transition.node.name))
		}
	}

	builder.WriteString("```\n")
	return builder.String()
}
